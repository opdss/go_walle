package version

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	//下面变量由编译时注入
	buildTimestamp  string
	buildCommitHash string
	buildVersion    string
	buildRelease    string

	Build build
)

func init() {
	Build = buildInfo()
}

// Info 构建版本信息
type build struct {
	Version    string    `json:"version"`   // tag版本号
	CommitHash string    `json:"hash"`      //commit hash
	Timestamp  time.Time `json:"timestamp"` //编译时间
	Release    bool      `json:"release"`   //是否发行版本
}

func (n build) String() string {
	return fmt.Sprintf("Release version：%s(%s);Release time：%s", n.Version, n.CommitHash, n.Timestamp)
}

func buildInfo() build {
	timestamp, err := strconv.ParseInt(buildTimestamp, 10, 64)
	if err != nil {
		timestamp = time.Now().Unix()
	}
	_info := build{
		Timestamp:  time.Unix(timestamp, 0),
		CommitHash: buildCommitHash,
		Version:    buildVersion,
		Release:    strings.ToLower(buildRelease) == "true",
	}
	return _info
}
