package model

import (
	"go-walle/app/model/field"
	"gorm.io/gorm"
	"time"
)

type Project struct {
	ID            int64 `gorm:"column:id;primaryKey;autoIncrement;" json:"id"`
	SpaceId       int64 `gorm:"column:space_id;index;not null;comment:所属空间" json:"space_id"`
	EnvironmentId int64 `gorm:"column:environment_id;not null;comment:所属环境" json:"environment_id"`
	UserId        int64 `gorm:"column:user_id;not null;comment:所属用户" json:"user_id"`

	Name         string `gorm:"column:name;size:100;not null;comment:名称" json:"name"`
	RepoUrl      string `gorm:"column:repo_url;size:500;not null;comment:仓库地址" json:"repo_url"`
	RepoMode     string `gorm:"column:repo_mode;size:10;not null;default:tag;comment:分支类型" json:"repo_mode"` // tag/branch
	RepoUsername string `gorm:"column:repo_username;size:100;not null;default:'';comment:仓库用户" json:"repo_username"`
	RepoPassword string `gorm:"column:repo_password;size:100;not null;default:'';comment:仓库密码" json:"repo_password"`
	RepoType     string `gorm:"column:repo_type;size:20;not null;default:git;comment:仓库类型" json:"repo_type"`

	Excludes  string `gorm:"column:excludes;size:1000;not null;default:'';comment:包含或者去除的文件列表" json:"excludes"` //包含或者去除的文件
	IsInclude int8   `gorm:"column:is_include;size:1;not null;default:0;comment:1去除0包含" json:"is_include"`      //是包含还是去除

	TaskVars    string `gorm:"column:task_vars;size:1000;not null;default:'';comment:全局环境变量" json:"task_vars"` //全局变量
	PrevDeploy  string `gorm:"column:prev_deploy;size:1000;not null;default:'';comment:编译前操作命令" json:"prev_deploy"`
	PostDeploy  string `gorm:"column:post_deploy;size:1000;not null;default:'';comment:编译后操作命令" json:"post_deploy"`
	PrevRelease string `gorm:"column:prev_release;size:1000;not null;default:'';comment:发布前操作命令" json:"prev_release"`
	PostRelease string `gorm:"column:post_release;size:1000;not null;default:'';comment:发布后操作命令" json:"post_release"`

	TargetRoot     string `gorm:"column:target_root;size:500;not null;default:'';comment:目标路径" json:"target_root"` //目标路径
	TargetReleases string `gorm:"column:target_releases;size:500;not null;default:'';comment:目标代码路径" json:"target_releases"`
	KeepVersionNum int    `gorm:"column:keep_version_num;size:4;not null;default:5;comment:保留版本数量" json:"keep_version_num"` //保留版本数量
	TaskAudit      int8   `gorm:"column:task_audit;size:1;not null;default:1;comment:上线单是否开启审核" json:"task_audit"`          //上线单是否开启审核
	Description    string `gorm:"column:description;size:500;not null;default:'';comment:简介说明" json:"description"`

	Master     string `gorm:"column:master" json:"master"`
	Version    string `gorm:"column:version;size:100;not null;default:'';comment:版本号" json:"version"`
	NoticeType string `gorm:"column:notice_type" json:"notice_type"`
	NoticeHook string `gorm:"column:notice_hook" json:"notice_hook"`

	Space       Space       `json:"space"`
	Environment Environment `json:"environment"`
	Servers     []Server    `gorm:"many2many:project_server" json:"servers"`

	Status field.Status `gorm:"column:status;size:1;not null;default:0;comment:状态" json:"status"`

	CreatedAt time.Time `gorm:"column:created_at;type:time;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:time;not null" json:"updated_at"`

	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}
