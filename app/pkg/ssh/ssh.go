package ssh

import (
	"fmt"
	"github.com/zeebo/errs"
	"golang.org/x/crypto/ssh"
	"os"
	"sync"
	"time"
)

var (
	ErrSSH = errs.Class("ssh")
)

type Ssh struct {
	config  *Config
	mux     *sync.Mutex
	clients map[string]*client
}

func NewSSH(conf *Config) (*Ssh, error) {
	return &Ssh{
		config:  conf,
		mux:     &sync.Mutex{},
		clients: make(map[string]*client),
	}, nil

	//go func() {
	//	tk := time.NewTicker(time.Second * 3)
	//	defer tk.Stop()
	//	for {
	//		select {
	//		case <-tk.C:
	//			if len(clients) == 0 {
	//				fmt.Printf("当前clients:[0] \r\n")
	//			} else {
	//				for _, v := range clients {
	//					fmt.Printf("当前clients:[%s], sessions:[%d]\r\n", v.serverConfig.String(), len(v.sessions))
	//				}
	//			}
	//		}
	//	}
	//}()
}

// Config ssh配置
type Config struct {
	IdentityFile     string        `help:"免密登陆密钥地址" default:"$HOME/.ssh/id_rsa"`
	IdentityPassword string        `help:"免密登陆密钥密码" default:""`
	Timeout          time.Duration `help:"连接超时" default:"30s"`
}

func (s *Ssh) IdentitySigners() (signers []ssh.Signer, err error) {
	var sg ssh.Signer
	sg, err = s.IdentitySigner()
	if err == nil {
		return []ssh.Signer{sg}, nil
	}
	return
}

func (s *Ssh) IdentitySigner() (signer ssh.Signer, err error) {
	_, err = os.Stat(s.config.IdentityFile)
	if err != nil {
		return nil, ErrSSH.New("ssh config IdentityFile: %s not exists", s.config.IdentityFile)
	}
	bytes, err := os.ReadFile(s.config.IdentityFile)
	if err != nil {
		return nil, ErrSSH.Wrap(err)
	}
	signer, err = ssh.ParsePrivateKey(bytes)
	if _, ok := err.(*ssh.PassphraseMissingError); ok {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(bytes, []byte(s.config.IdentityPassword))
	}
	if err != nil {
		err = ErrSSH.Wrap(err)
	}
	return
}

type ServerConfig struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"` //如果密码为空，则认为是免密登陆
	Port     int    `json:"port"`
}

func (s *ServerConfig) String() (key string) {
	return fmt.Sprintf("%s:%s@%s:%d", s.User, s.Password, s.Host, s.Port)
}

func (s *Ssh) NewClient(conf ServerConfig) (*ssh.Client, error) {
	sc, err := newSshClient(s, &conf)
	if err != nil {
		return nil, err
	}
	return sc.client, nil
}

func (s *Ssh) newClient(conf ServerConfig) (sc *client, err error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	key := conf.String()
	if v, ok := s.clients[key]; ok {
		return v, nil
	}
	sc, err = newSshClient(s, &conf)
	if err != nil {
		return
	}
	s.clients[key] = sc
	return
}

func (s *Ssh) removeClient(conf *ServerConfig) {
	s.mux.Lock()
	defer s.mux.Unlock()
	key := conf.String()
	if _, ok := s.clients[key]; ok {
		delete(s.clients, key)
	}
}

// NewTerminal 获取会话终端
func (s *Ssh) NewTerminal(conf ServerConfig, cols, rows int) (sess *Terminal, err error) {
	sshClient, err := s.newClient(conf)
	if err != nil {
		err = ErrSSH.Wrap(err)
		return
	}
	return sshClient.newTerminal(cols, rows)
}

// RunCmd 直接连接执行命令
func (s *Ssh) RunCmd(conf ServerConfig, cmd string) (output []byte, err error) {
	sshClient, err := s.newClient(conf)
	if err != nil {
		err = ErrSSH.Wrap(err)
		return
	}
	return sshClient.RunCmd(cmd)
}

func (s *Ssh) NewSftp(conf ServerConfig) (*Sftp, error) {
	sshClient, err := s.newClient(conf)
	if err != nil {
		err = ErrSSH.Wrap(err)
		return nil, err
	}
	return sshClient.newSftp()
}

func (s *Ssh) NewRemoteExec(conf ServerConfig) (*RemoteExec, error) {
	sshClient, err := s.newClient(conf)
	if err != nil {
		err = ErrSSH.Wrap(err)
		return nil, err
	}
	return sshClient.newRemoteExec()
}
