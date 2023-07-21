package ctx

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-walle/app/global"
	"go-walle/app/pkg/jwt"
	"go-walle/app/service/common"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const SpaceHeaderName = "Space-Id"
const AuthCtxKey = "__AUTH_KEY__"
const RoleCtxKey = "__ROLE_KEY__"

func GetSpaceId(ctx *gin.Context) int64 {
	n, err := strconv.Atoi(strings.Trim(ctx.GetHeader(SpaceHeaderName), " "))
	if err != nil {
		return 0
	}
	return int64(n)
}

func GetBearerToken(ctx *gin.Context) string {
	tokenStr := ctx.Request.Header.Get("Authorization")
	if len(tokenStr) > 0 && tokenStr[0:7] == "Bearer " {
		return tokenStr[7:]
	}
	//websocket
	tokenStr = ctx.Request.Header.Get("Sec-WebSocket-Protocol")
	res := strings.Split(tokenStr, ",")
	if len(res) == 2 {
		tokenStr = res[0]
		ctx.Request.Header.Set(SpaceHeaderName, res[1])
	}
	return tokenStr
}

func ValidateBearerToken(ctx *gin.Context) (*jwt.TokenClaims, error) {
	token := GetBearerToken(ctx)
	if len(token) == 0 {
		return nil, errors.New("bearer token not exists")
	}
	return global.Jwt.ValidateToken(token)
}

func IsLogin(ctx *gin.Context) (auth *jwt.TokenClaims, err error) {
	if v, ok := ctx.Value(AuthCtxKey).(*jwt.TokenClaims); ok {
		return v, nil
	}
	auth, err = ValidateBearerToken(ctx)
	if err != nil {
		return
	}
	ctx.Set(AuthCtxKey, auth)
	return
}

// UserId 当前登陆用户id
func UserId(ctx *gin.Context) int64 {
	auth, err := IsLogin(ctx)
	if err != nil {
		return 0
	}
	return auth.UserId
}

// GetSpaceWithId 当前登陆用户id
func GetSpaceWithId(ctx *gin.Context) (*common.SpaceWithId, error) {
	id := ctx.Param("id")
	n, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.New("not found id")
	}
	return &common.SpaceWithId{SpaceId: GetSpaceId(ctx), ID: int64(n)}, nil
}

// Username 当前登陆用户id
func Username(ctx *gin.Context) string {
	auth, err := IsLogin(ctx)
	if err != nil {
		return ""
	}
	return auth.Username
}

// SetRole 当前登陆用户id
func SetRole(ctx *gin.Context, role string) {
	ctx.Set(RoleCtxKey, role)
}

func UpGrader(ctx *gin.Context) (*websocket.Conn, error) {
	upGrader := websocket.Upgrader{
		HandshakeTimeout: 10 * time.Second,
		// cross origin domain
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		// 处理 Sec-WebSocket-Protocol Header
		Subprotocols: []string{GetBearerToken(ctx), strconv.Itoa(int(GetSpaceId(ctx)))},
	}
	return upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
}
