package project

import "go-walle/app/pkg/db"

type CreateReq struct {
	SpaceId       int64  `json:"-" binding:"required,gt=0"`
	Name          string `json:"name" binding:"required,max=50"`
	EnvironmentId int64  `json:"environment_id" binding:"required,gt=0"`
	RepoUrl       string `json:"repo_url" binding:"required,max=500"`
	RepoType      string `json:"repo_type" binding:"required,max=20"`
	RepoMode      string `json:"repo_mode" binding:"required,max=20"`

	ServerIds      []int64 `json:"server_ids" binding:"required,unique,dive,gt=0"`
	TargetRoot     string  `json:"target_root" binding:"required,max=100"`
	TargetReleases string  `json:"target_releases" binding:"required,max=100"`
	KeepVersionNum int     `json:"keep_version_num" binding:"required,gt=0"`

	Excludes    string `json:"excludes" binding:"omitempty"`
	IsInclude   int8   `json:"is_include" binding:"omitempty"`
	TaskVars    string `json:"task_vars" binding:"omitempty"`
	PrevDeploy  string `json:"prev_deploy" binding:"omitempty"`
	PostDeploy  string `json:"post_deploy" binding:"omitempty"`
	PrevRelease string `json:"prev_release" binding:"omitempty"`
	PostRelease string `json:"post_release" binding:"omitempty"`

	TaskAudit int8 `json:"task_audit" binding:"omitempty"`

	Description string `json:"description" binding:"omitempty,max=500"`
}

type UpdateReq struct {
	SpaceId int64 `json:"-" binding:"required,gt=0"`

	ID            int64  `json:"id" validate:"required,gt=0"`
	Name          string `json:"name" binding:"required,max=50"`
	EnvironmentId int64  `json:"environment_id" binding:"required,gt=0"`
	RepoUrl       string `json:"repo_url" binding:"required,max=500"`
	RepoType      string `json:"repo_type" binding:"required,max=20"`
	RepoMode      string `json:"repo_mode" binding:"required,max=20"`

	ServerIds      []int64 `json:"server_ids" binding:"required,unique,dive,gt=0"`
	TargetRoot     string  `json:"target_root" binding:"required,max=100"`
	TargetReleases string  `json:"target_releases" binding:"required,max=100"`
	KeepVersionNum int     `json:"keep_version_num" binding:"required,gt=0"`

	Excludes    string `json:"excludes" binding:"omitempty"`
	IsInclude   int8   `json:"is_include" binding:"omitempty"`
	TaskVars    string `json:"task_vars" binding:"omitempty"`
	PrevDeploy  string `json:"prev_deploy" binding:"omitempty"`
	PostDeploy  string `json:"post_deploy" binding:"omitempty"`
	PrevRelease string `json:"prev_release" binding:"omitempty"`
	PostRelease string `json:"post_release" binding:"omitempty"`

	TaskAudit int8 `json:"task_audit" binding:"omitempty"`

	Description string `json:"description" binding:"omitempty,max=500"`
}

func (r *UpdateReq) Fields() []string {
	return []string{
		"name", "environment_id", "repo_url", "repo_type", "repo_mode",
		"target_root", "target_releases", "keep_version_num",
		"excludes", "is_include", "task_vars", "prev_deploy", "post_deploy", "prev_release", "post_release",
		"task_audit", "description",
	}
}

type ListReq struct {
	SpaceId int64 `json:"-" binding:"required,gt=0"`

	EnvironmentId int64 `json:"environment_id" query:"environment_id" form:"environment_id" binding:"omitempty,gt=0"`
	//WithLastRelease bool  `json:"with_last_release" query:"environment_id" form:"environment_id" binding:"omitempty"`
	db.Paginator
}

type DetectionMsg struct {
	Title string `json:"title"`
	Error string `json:"error"`
	Todo  string `json:"todo"`
}
