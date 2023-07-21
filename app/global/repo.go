package global

import "go-walle/app/pkg/repo"

var Repo *repo.Repos

func initRepo(conf *repo.Config) (err error) {
	Repo = repo.NewRepos(conf)
	return
}
