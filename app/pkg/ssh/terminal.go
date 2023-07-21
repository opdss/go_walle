package ssh

import (
	"golang.org/x/crypto/ssh"
	"io"
)

// Terminal  模拟终端会话
type Terminal struct {
	sid     string
	client  *client
	session *ssh.Session
	reader  io.Reader
	writer  io.Writer
}

// Close 关闭终端会话
func (s *Terminal) Close() error {
	defer func() {
		s.client.done()
	}()
	return s.session.Close()
}

// Read 返回数据
func (s *Terminal) Read(p []byte) (n int, err error) {
	return s.reader.Read(p)
}

// Write 写入数据
func (s *Terminal) Write(p []byte) (n int, err error) {
	return s.writer.Write(p)
}

// WindowChange 改变窗口大小
func (s *Terminal) WindowChange(h, w int) error {
	return s.session.WindowChange(h, w)
}
