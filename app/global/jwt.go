package global

import "go-walle/app/pkg/jwt"

var Jwt *jwt.Jwt

func initJwt(conf *jwt.Config) (err error) {
	Jwt, err = jwt.NewJWT(conf)
	return
}
