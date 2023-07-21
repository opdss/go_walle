package jwt

import (
	jwtgo "github.com/form3tech-oss/jwt-go"
	"github.com/zeebo/errs"
	"time"
)

var Error = errs.Class("JWT")

type Jwt struct {
	*Config
}

type Config struct {
	TokenExpire        time.Duration `help:"token有效期" devDefault:"24h0m0s" default:"30m0s"`
	RefreshTokenExpire time.Duration `help:"刷新token有效期" default:"720h0m0s"`
	JWTKey             string        `help:"JWT加密key" default:"|^_^|"`
	JWTAlg             string        `help:"JWT加密方式" default:"HS256"`
}

func NewJWT(cfg *Config) (*Jwt, error) {
	return &Jwt{cfg}, nil
}

type TokenClaims struct {
	*jwtgo.StandardClaims
	TokenPayload
}

type TokenPayload struct {
	UserId    int64  `json:"uid"`
	Username  string `json:"une"`
	Email     string `json:"eml"`
	IsRefresh bool   `json:"irf"`
}

// CreateToken 创建jwt token
func (j *Jwt) CreateToken(ap TokenPayload) (string, int64, error) {
	return j.createToken(ap, j.TokenExpire)
}

// CreateRefreshToken 创建jwt token
func (j *Jwt) CreateRefreshToken(ap TokenPayload) (string, int64, error) {
	return j.createToken(ap, j.RefreshTokenExpire)
}

// ValidateToken 校验jwt token
func (j *Jwt) ValidateToken(token string) (*TokenClaims, error) {
	parser := &jwtgo.Parser{}
	tk, err := parser.ParseWithClaims(token, &TokenClaims{}, func(token *jwtgo.Token) (interface{}, error) {
		return []byte(j.JWTKey), nil
	})
	if err != nil {
		return nil, Error.Wrap(err)
	}
	if c, ok := tk.Claims.(*TokenClaims); ok {
		return c, nil
	}
	return nil, Error.New("claims parse error")
}

func (j *Jwt) createToken(ap TokenPayload, exp time.Duration) (string, int64, error) {
	t := jwtgo.New(jwtgo.GetSigningMethod(j.JWTAlg))
	expire := time.Now().Add(exp).Unix()
	t.Claims = &TokenClaims{
		&jwtgo.StandardClaims{
			ExpiresAt: expire,
		},
		ap,
	}
	a, err := t.SignedString([]byte(j.JWTKey))
	return a, expire, Error.Wrap(err)
}
