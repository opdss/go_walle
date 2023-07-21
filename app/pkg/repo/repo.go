package repo

import (
	"github.com/zeebo/errs"
	"time"
)

type TypeRepo string

const (
	GitRepo TypeRepo = "git"
	SvnRepo TypeRepo = "svn"
)

var (
	ErrRepo    = errs.Class("repo")
	ErrRepoGit = errs.Class("repo.git")
	ErrRepoSvn = errs.Class("repo.svn")
)

type Repos struct {
	config *Config
}

type Config struct {
	RepoDir string `help:"代码本地存放目录" devDefault:"$ROOT/warehouse" default:"/var/lib/walle/warehouse"`
	Git     GitConfig
	Svn     SvnConfig
}

func NewRepos(cfg *Config) *Repos {
	return &Repos{
		config: cfg,
	}
}

type Tag struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
}

type Commit struct {
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Hash      string    `json:"hash"`
}

type Branch struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
}

type Repo interface {
	Tags() ([]Tag, error)
	Commits(branch string) ([]Commit, error)
	Branches() ([]Branch, error)
	CheckoutToBranch(branch string) error
	CheckoutToCommit(branch, commit string) error
	CheckoutToTag(tag string) error
	Path() string
	Type() TypeRepo
}

func (r *Repos) New(repoType TypeRepo, repoUrl, projectName string) (Repo, error) {
	switch repoType {
	case GitRepo:
		return NewGit(&r.config.Git, repoUrl, r.config.RepoDir+"/"+projectName)
	case SvnRepo:
		return NewSvn(&r.config.Svn, repoUrl, r.config.RepoDir+"/"+projectName)
	}
	return nil, ErrRepo.New("仓库类型不支持")
}
