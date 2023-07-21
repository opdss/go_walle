package notice

type Config struct {
	Dingtalk DingtalkConfig
	Email    EmailConfig
}

type Notice interface {
	Send() error
}
