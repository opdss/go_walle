package ssh

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"sync"
)

// client ssh方便复用tcp连接管理
type client struct {
	mux          sync.Mutex
	sh           *Ssh
	serverConfig *ServerConfig
	client       *ssh.Client
	ref          int
}

func (s *client) Key() string {
	return s.serverConfig.String()
}

func (s *client) add() {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.ref++
}

func (s *client) done() {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.ref--
	if s.ref == 0 {
		_ = s.client.Close()
		s.sh.removeClient(s.serverConfig)
	}
}

func newSshClient(sh *Ssh, conf *ServerConfig) (_ *client, err error) {
	config := &ssh.ClientConfig{
		User:            conf.User,
		Auth:            []ssh.AuthMethod{ssh.Password(conf.Password), ssh.PublicKeysCallback(sh.IdentitySigners)},
		Timeout:         sh.config.Timeout,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	config.SetDefaults()
	tcpAddress := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	sshClient, err := ssh.Dial("tcp", tcpAddress, config)
	if nil != err {
		return nil, err
	}
	return &client{
		serverConfig: conf,
		client:       sshClient,
	}, nil
}

func (s *client) RunCmd(cmd string) (output []byte, err error) {
	var session *ssh.Session
	session, err = s.client.NewSession()
	if err != nil {
		return nil, err
	}
	s.mux.Lock()
	s.add()
	s.mux.Unlock()
	output, err = session.CombinedOutput(cmd)
	_ = session.Close()
	s.mux.Lock()
	s.done()
	s.mux.Unlock()
	return
}

func (s *client) newTerminal(cols, rows int) (term *Terminal, err error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	session, err := s.client.NewSession()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = session.Close()
		}
	}()
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	err = session.RequestPty("xterm", rows, cols, modes)
	if err != nil {
		return
	}
	var reader io.Reader
	var writer io.Writer
	reader, err = session.StdoutPipe()
	if err != nil {
		fmt.Println("session.StdoutPipe error:", err)
		return
	}
	writer, err = session.StdinPipe()
	if err != nil {
		fmt.Println("session.StdinPipe error", err)
		return
	}
	err = session.Shell()
	if err != nil {
		return
	}
	term = &Terminal{
		client:  s,
		session: session,
		reader:  reader,
		writer:  writer,
	}
	s.add()
	return term, nil
}

func (s *client) newSftp() (*Sftp, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	scp, err := sftp.NewClient(s.client)
	if err != nil {
		return nil, err
	}
	_sftp := &Sftp{
		client:     s,
		sftpClient: scp,
	}
	s.add()
	return _sftp, nil
}

func (s *client) newRemoteExec() (*RemoteExec, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	_re := &RemoteExec{
		client: s,
		envs:   NewEnvs(),
	}
	s.add()
	return _re, nil
}
