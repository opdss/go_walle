package errcode

import "fmt"

type ErrCode int

const (
	Successes        ErrCode = 0  //ok
	ErrServer        ErrCode = 1  //服务器处理失败
	ErrInvalidParams ErrCode = 2  //无效参数
	ErrTimeOut       ErrCode = 3  //处理超时
	ErrForbidden     ErrCode = 4  //无权限访问
	ErrUnauthorized  ErrCode = 5  //未授权
	ErrRequest       ErrCode = 6  //请求错误
	ErrInvalidPwd    ErrCode = 7  //用户名或密码错误
	ErrUserDisabled  ErrCode = 8  //账号被禁用
	ErrCaptcha       ErrCode = 9  //验证码错误
	ErrNotFound      ErrCode = 10 //资源未找到或没有权限访问
	ErrDataHasExist  ErrCode = 11 //数据已经存在
)

type ErrWrap struct {
	errCode ErrCode
	err     error
}

func (i ErrCode) Error() string {
	return i.String()
}

func (i ErrCode) New(format string, args ...any) ErrWrap {
	return ErrWrap{
		errCode: i,
		err:     fmt.Errorf(format, args...),
	}
}

func (i ErrCode) Wrap(err error) ErrWrap {
	return ErrWrap{errCode: i, err: err}
}

func (ew ErrWrap) Error() string {
	return ew.err.Error()
}

func (ew ErrWrap) ErrCode() int {
	return int(ew.errCode)
}
