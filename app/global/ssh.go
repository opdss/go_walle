package global

import (
	"go-walle/app/pkg/ssh"
)

var Ssh *ssh.Ssh

func initSsh(conf *ssh.Config) (err error) {
	Ssh, err = ssh.NewSSH(conf)
	return
}
