package repo

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	gitConfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"os"
	"strings"
)

type GitConfig struct {
	Username           string `help:"BasicAuth授权模式下，git帐号" default:"git"`
	Password           string `help:"BasicAuth授权模式下，git密码" devDefault:"C-Bp2FAS2yBzYCuE7ECj" default:""`
	PrivateKeyFile     string `help:"免密模式下，私钥证书地址" default:"$HOME/.ssh/id_rsa"`
	PrivateKeyUsername string `help:"免密模式下，私钥证书用户" default:"git"`
	PrivateKeyPassword string `help:"免密模式下，私钥证书密码" default:""`
	FetchDepth         int    `help:"记录数量" default:"50"`
}

type Git struct {
	config  *GitConfig
	path    string
	repoUrl string
	repo    *git.Repository

	auth transport.AuthMethod //auth
}

func NewGit(cfg *GitConfig, url string, path string) (*Git, error) {
	repo := &Git{
		config:  cfg,
		path:    path,
		repoUrl: url,
	}
	r, err := repo.init()
	if err != nil {
		err = ErrRepoGit.Wrap(err)
	}
	return r, err
}

func (srv *Git) init() (_ *Git, err error) {
	//设置auth
	srv.auth, err = srv.getAuth()
	if err != nil {
		return
	}
	//读取存在的目录，项目不存在，则重新clone
	srv.repo, err = git.PlainOpen(srv.path)
	if err != nil {
		if errors.Is(err, git.ErrRepositoryNotExists) {
			err = srv.clone()
			if err == nil {
				return srv, nil
			}
		}
		return
	}
	//检查存在项目是否跟远程地址一致
	remote, _ := srv.repo.Remote(git.DefaultRemoteName)
	if strings.Contains(remote.String(), srv.repoUrl) {
		return srv, nil
	}
	return nil, fmt.Errorf("dir[%s] remote:%s, not:%s", srv.path, remote.String(), srv.repoUrl)
}

// clone 克隆项目到目录， 如果存在目录，则删除
func (srv *Git) clone() (err error) {
	_ = os.RemoveAll(srv.path)
	srv.repo, err = git.PlainClone(srv.path, false, &git.CloneOptions{
		Auth:     srv.auth,
		URL:      srv.repoUrl,
		Progress: os.Stdout,
	})
	return err
}

func (srv *Git) pull(branch string) error {
	w, err := srv.repo.Worktree()
	if err != nil {
		return err
	}
	err = w.Pull(&git.PullOptions{
		Auth:          srv.auth,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
		Progress:      os.Stdout,
	})
	if errors.Is(err, git.NoErrAlreadyUpToDate) {
		return nil
	}
	if errors.Is(err, plumbing.ErrReferenceNotFound) {
		_ = srv.delBranch(branch)
		_ = srv.repo.Storer.RemoveReference(plumbing.NewBranchReferenceName(branch))
	}
	return err
}

func (srv *Git) fetch() error {
	err := srv.repo.Fetch(&git.FetchOptions{
		RefSpecs: []gitConfig.RefSpec{"refs/*:refs/*", "HEAD:ref/heads/HEAD"},
		Auth:     srv.auth,
	})
	//r, _ := srv.repo.References()
	//for {
	//	s, err := r.Next()
	//	if err != nil {
	//		if err == io.EOF {
	//			break
	//		}
	//		return err
	//	}
	//	fmt.Println(s.Name(), s.Type(), s.Hash())
	//}
	if errors.Is(err, git.NoErrAlreadyUpToDate) {
		return nil
	}
	return err
}

func (srv *Git) CheckoutToBranch(branch string) (err error) {
	defer func() {
		if err != nil {
			err = ErrRepoGit.Wrap(err)
		}
	}()
	err = srv.fetch()
	if err != nil {
		return
	}
	var w *git.Worktree
	w, err = srv.repo.Worktree()
	if err != nil {
		return
	}
	return w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
		Force:  true,
	})
	//fmt.Println(444, err)
	//if err != nil {
	//	return
	//}
	//return
	//err = srv.Pull(branch)
	//fmt.Println(555, err)
	//return
}

func (srv *Git) CheckoutToCommit(branch, commit string) (err error) {
	defer func() {
		if err != nil {
			err = ErrRepoGit.Wrap(err)
		}
	}()
	err = srv.CheckoutToBranch(branch)
	if err != nil {
		return
	}
	var w *git.Worktree
	w, err = srv.repo.Worktree()
	if err != nil {
		return
	}
	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(commit),
	})
	return
}

func (srv *Git) CheckoutToTag(tag string) (err error) {
	defer func() {
		if err != nil {
			err = ErrRepoGit.Wrap(err)
		}
	}()
	err = srv.fetch()
	if err != nil {
		return
	}
	var w *git.Worktree
	w, err = srv.repo.Worktree()
	if err != nil {
		return
	}
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewTagReferenceName(tag),
	})
	return
}

func (srv *Git) getBranch(branch string) (b *gitConfig.Branch, err error) {
	b, err = srv.repo.Branch(branch)
	if err == nil {
		return
	}
	b = &gitConfig.Branch{
		Name:   branch,
		Merge:  plumbing.NewBranchReferenceName(branch),
		Rebase: "true",
	}
	err = srv.repo.CreateBranch(b)
	return
}

func (srv *Git) delBranch(branch string) (err error) {
	_ = srv.repo.DeleteBranch(branch)
	return
}

// Branches 获取所有远程分支
func (srv *Git) Branches() ([]Branch, error) {
	_branches := make([]Branch, 0)
	remote, err := srv.repo.Remote(git.DefaultRemoteName)
	if err != nil {
		return nil, ErrRepoGit.Wrap(err)
	}
	br, err := remote.List(&git.ListOptions{Auth: srv.auth})
	if err != nil {
		return nil, ErrRepoGit.Wrap(err)
	}
	for _, v := range br {
		if v.Name().IsBranch() {
			_Branch := Branch{
				Name: v.Name().Short(),
				Hash: v.Hash().String(),
			}
			_branches = append(_branches, _Branch)
		}
	}
	return _branches, nil
}

// Tags 获取所有标签
func (srv *Git) Tags() ([]Tag, error) {
	_tags := make([]Tag, 0)
	remote, err := srv.repo.Remote(git.DefaultRemoteName)
	if err != nil {
		return nil, ErrRepoGit.Wrap(err)
	}
	br, err := remote.List(&git.ListOptions{Auth: srv.auth})
	if err != nil {
		return nil, ErrRepoGit.Wrap(err)
	}
	for _, v := range br {
		if v.Name().IsTag() {
			_tag := Tag{
				Name: v.Name().Short(),
				Hash: v.Hash().String(),
			}
			_tags = append(_tags, _tag)
		}
	}
	return _tags, nil
}

// Commits 获取对应的分支的commits
func (srv *Git) Commits(branch string) ([]Commit, error) {
	r, err := srv.repo.Reference(plumbing.NewRemoteReferenceName(git.DefaultRemoteName, branch), false)
	if err != nil {
		return nil, ErrRepoGit.Wrap(err)
	}
	_commits := make([]Commit, 0)
	br, err := srv.repo.Log(&git.LogOptions{From: r.Hash(), Order: git.LogOrderCommitterTime})
	if err != nil {
		return nil, ErrRepoGit.Wrap(err)
	}
	_ = br.ForEach(func(commit *object.Commit) error {
		_commits = append(_commits, Commit{
			Name:      commit.Hash.String()[:8] + "#" + commit.Message,
			Message:   commit.Message,
			Timestamp: commit.Committer.When,
			Hash:      commit.Hash.String(),
		})
		return nil
	})
	if srv.config.FetchDepth > 0 && len(_commits) > srv.config.FetchDepth {
		_commits = _commits[:srv.config.FetchDepth]
	}
	return _commits, nil
}

func (srv *Git) getAuth() (auth transport.AuthMethod, _ error) {
	if srv.repoUrl[0:3] == "git" {
		_, err := os.Stat(srv.config.PrivateKeyFile)
		if err != nil {
			return nil, fmt.Errorf("git config PrivateKeyFile: %s not exisit\n", srv.config.PrivateKeyFile)
		}
		publicKeys, err := ssh.NewPublicKeysFromFile(srv.config.PrivateKeyUsername, srv.config.PrivateKeyFile, srv.config.PrivateKeyPassword)
		if err != nil {
			return nil, err
		}
		return publicKeys, nil
	}
	return &http.BasicAuth{
		Username: srv.config.Username,
		Password: srv.config.Password,
	}, nil
}

func (srv *Git) Path() string {
	return srv.path
}

func (srv *Git) Type() TypeRepo {
	return GitRepo
}
