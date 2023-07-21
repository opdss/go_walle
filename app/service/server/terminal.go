package server

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"go-walle/app/model"
	"go-walle/app/model/field"
	"go-walle/app/pkg/ssh"
	"go-walle/app/service/common"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
	"time"
	"unicode/utf8"
)

const (
	waringMsg = iota
	errorMsg
	successMsg
	defaultMsg
)

const (
	connectTimeout     = time.Minute * 10 //保持连接最长时间
	buffTime           = time.Microsecond * 500
	wsMsgTypeResize    = "resize"
	wsMsgTypeCmd       = "cmd"
	wsMsgTypeHeartbeat = "ping"
)

type TerminalWsMsg struct {
	Typ string `json:"typ"`
	Cmd string `json:"cmd,omitempty"`
	Col int    `json:"col,omitempty"`
	Row int    `json:"row,omitempty"`
}

func (srv *Service) Check(spaceWithId *common.SpaceWithId) error {
	serverDetail := model.Server{}
	err := srv.db.Where(spaceWithId).First(&serverDetail).Error
	if err != nil {
		return err
	}
	output, err := srv.ssh.RunCmd(ssh.ServerConfig{
		User: serverDetail.User,
		Host: serverDetail.Host,
		Port: serverDetail.Port,
	}, "pwd")
	srv.log.Debug("CheckConnect", zap.String("cmd", "pwd"), zap.ByteString("output", output), zap.Error(err))
	if err != nil && serverDetail.Status.IsEnable() {
		return srv.db.Model(&serverDetail).Where("id=?", serverDetail.ID).UpdateColumn("status", field.StatusDisable).Error
	}
	if err == nil && serverDetail.Status.IsDisable() {
		return srv.db.Model(&serverDetail).Where("id=?", serverDetail.ID).UpdateColumn("status", field.StatusEnable).Error
	}
	return err
}

func (srv *Service) SetAuthorized(params *SetAuthorizedReq) error {
	serverDetail := model.Server{SpaceId: params.SpaceId, ID: params.ID}
	err := srv.db.Where(serverDetail).First(&serverDetail).Error
	if err != nil {
		return err
	}
	if serverDetail.Status.IsEnable() {
		return errors.New("该服务器能正常连接，无需设置")
	}
	signer, err := srv.ssh.IdentitySigner()
	if err != nil {
		return err
	}
	hostname, _ := os.Hostname()
	publicKeyStr := fmt.Sprintf("%s %s %s", signer.PublicKey().Type(), base64.StdEncoding.EncodeToString(signer.PublicKey().Marshal()), hostname)
	runCmd := fmt.Sprintf("mkdir -p $HOME/.ssh && echo '%s' >> $HOME/.ssh/authorized_keys && chmod 600 $HOME/.ssh/authorized_keys", publicKeyStr)
	output, err := srv.ssh.RunCmd(ssh.ServerConfig{
		User:     serverDetail.User,
		Host:     serverDetail.Host,
		Password: params.Password,
		Port:     serverDetail.Port,
	}, runCmd)
	srv.log.Debug("Setting", zap.String("cmd", runCmd), zap.ByteString("output", output), zap.Error(err))
	if err == nil {
		_err := srv.db.Model(&serverDetail).Where("id=?", serverDetail.ID).UpdateColumn("status", field.StatusEnable).Error
		if _err != nil {
			srv.log.Error("更新数据库失败", zap.Int64("server_id", serverDetail.ID), zap.Int("status", field.StatusEnable))
		}
	}
	return err
}

func (srv *Service) Terminal(wsConn *websocket.Conn, spaceWithId *common.SpaceWithId, username string) error {
	wsSendMsg := func(msg string, msgType int) error {
		_err := wsConn.WriteMessage(websocket.TextMessage, []byte(terminalMsg(msg, msgType)))
		if _err != nil {
			srv.log.Debug("发送ws数据失败", zap.Error(_err))
		}
		return _err
	}
	serverDetail := model.Server{}
	err := srv.db.Where(spaceWithId).First(&serverDetail).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = wsSendMsg("该服务器不存在！", errorMsg)
		} else {
			_ = wsSendMsg(err.Error(), errorMsg)
		}
		return err
	}

	if err = wsSendMsg("正在连接服务器...", successMsg); err != nil {
		return err
	}
	sshTerminal, err := srv.ssh.NewTerminal(ssh.ServerConfig{
		User:     serverDetail.User,
		Host:     serverDetail.Host,
		Password: "",
		Port:     serverDetail.Port,
	}, 200, 40)
	if err != nil {
		_ = wsSendMsg(err.Error(), errorMsg)
		return err
	}
	defer func() {
		_ = sshTerminal.Close()
	}()
	if err = wsSendMsg("连接服务器成功！", successMsg); err != nil {
		return err
	}
	if err = wsSendMsg("Hello "+username+"，您所操作的所有命令都将会被记录，请谨慎操作！！！", waringMsg); err != nil {
		return err
	}
	srv.dealMsg(wsConn, sshTerminal)
	return nil
}

// dealMsg 终端数据交互
func (srv *Service) dealMsg(wsConn *websocket.Conn, sshTerminal *ssh.Terminal) {
	connectTimeoutT := time.NewTimer(connectTimeout)
	bufTimeT := time.NewTimer(buffTime)
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		connectTimeoutT.Stop()
		bufTimeT.Stop()
	}()

	go func() {
		var err error
		defer func() {
			if err != nil {
				srv.log.Error("处理ws消息出错", zap.Error(err))
			}
		}()
		var msg []byte
		for {
			select {
			//监听上下文退出
			case <-ctx.Done():
				return
			default:
				_, msg, err = wsConn.ReadMessage()
				if err != nil {
					cancel()
					return
				}
				wsMsg := new(TerminalWsMsg)
				if err = json.Unmarshal(msg, wsMsg); err != nil {
					continue
				}
				switch wsMsg.Typ {
				case wsMsgTypeResize:
					err = sshTerminal.WindowChange(wsMsg.Row, wsMsg.Col)
				case wsMsgTypeHeartbeat:
				default:
					_, err = sshTerminal.Write([]byte(wsMsg.Cmd))
				}
				if err != nil {
					cancel()
					return
				}
			}
		}
	}()

	r := make(chan rune)
	//读取buf
	br := bufio.NewReader(sshTerminal)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				x, size, err := br.ReadRune()
				if err != nil {
					cancel()
					srv.log.Error("处理shell消息出错", zap.Error(err))
					return
				}
				if size > 0 {
					r <- x
				}
			}
		}
	}()

	buf := make([]byte, 0)
	// 主循环
	for {
		select {
		case <-connectTimeoutT.C:
			cancel()
			return
		case <-ctx.Done():
			return
		case <-bufTimeT.C:
			if len(buf) != 0 {
				err := wsConn.WriteMessage(websocket.TextMessage, buf)
				buf = []byte{}
				if err != nil {
					cancel()
					srv.log.Error("ws.WriteMessage:", zap.Error(err))
					return
				}
				connectTimeoutT.Reset(connectTimeout)
			}
			bufTimeT.Reset(buffTime)
		case d := <-r:
			if d != utf8.RuneError {
				p := make([]byte, utf8.RuneLen(d))
				utf8.EncodeRune(p, d)
				buf = append(buf, p...)
			} else {
				buf = append(buf, []byte("@")...)
			}
			connectTimeoutT.Reset(connectTimeout)
		}
	}
}

func terminalMsg(msg string, typ int) string {
	switch typ {
	case waringMsg:
		return fmt.Sprintf("\x1b[33m%s\x1b[m\r\n", msg)
	case errorMsg:
		return fmt.Sprintf("\x1b[31m%s\x1b[m\r\n", msg)
	case successMsg:
		return fmt.Sprintf("\x1b[32m%s\x1b[m\r\n", msg)
	default:
		return fmt.Sprintf("%s\r\n", msg)
	}
}
