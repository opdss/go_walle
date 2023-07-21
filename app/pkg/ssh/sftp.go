package ssh

import (
	"github.com/pkg/sftp"
	"io"
	"os"
)

type Sftp struct {
	client     *client
	sftpClient *sftp.Client
}

// Close 关闭终端会话
func (s *Sftp) Close() error {
	defer func() {
		s.client.done()
	}()
	return s.sftpClient.Close()
}

func (s *Sftp) Copy(localFile, remoteFile string) error {
	lf, err := os.Open(localFile)
	if err != nil {
		return err
	}
	rf, err := s.sftpClient.Create(remoteFile)
	if err != nil {
		return err
	}
	_, err = io.Copy(rf, lf)
	return err
}
