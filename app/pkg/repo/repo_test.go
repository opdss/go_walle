package repo

import (
	"fmt"
	"testing"
)

var dir = "/tmp/test_git"

//var url = "http://code.nextstorage.cn/storj/punch.git"

var testConfig Config = Config{
	RepoDir: "/Users/wuxin/worker/go-walle/warehouse",
	Git: GitConfig{
		Username:           "git",
		Password:           "C-Bp2FAS2yBzYCuE7ECj",
		PrivateKeyFile:     "/Users/wuxin/.ssh/id_rsa",
		PrivateKeyUsername: "git",
		PrivateKeyPassword: "",
	},
}

var url = "git@code.nextstorage.cn:wuxin/tools.git"
var name = "itools"

func TestRepo(t *testing.T) {
	_ = Init(&testConfig)
	repo, err := New(RepoGit, url, name)
	fmt.Println(repo, err)
}

//
//func TestTags(t *testing.T) {
//	_ = Init(&testConfig)
//	repo, err := New(RepoType("git"), url, "itools")
//	if err != nil {
//		t.Error(err)
//	}
//	tags, err := repo.Tags()
//	if err != nil {
//		t.Error(err)
//	}
//	fmt.Println("tags", tags)
//}
//
//func TestBranches(t *testing.T) {
//	_ = Init(&testConfig)
//	repo, err := New(RepoType("git"), url, "itools")
//	if err != nil {
//		t.Error(err)
//	}
//	branches, err := repo.Branches()
//	if err != nil {
//		t.Error(err)
//	}
//	fmt.Println("branches", branches)
//}
//
//func TestCommits(t *testing.T) {
//	_ = Init(&testConfig)
//	repo, err := New(RepoType("git"), url, "itools")
//	if err != nil {
//		t.Error(err)
//	}
//	tags, err := repo.Commits("main")
//	if err != nil {
//		t.Error(err)
//	}
//	fmt.Println("commits", tags)
//}

//func TestCheckoutToBranch(t *testing.T) {
//	_ = Init(&testConfig)
//	repo, err := New(RepoType("git"), url, "itools")
//	if err != nil {
//		t.Error(err)
//	}
//	err = repo.CheckoutToBranch("test")
//	if err != nil {
//		t.Error(err)
//	}
//	fmt.Println("CheckoutToBranch success", repo.Path())
//}

//func TestCheckoutToCommit(t *testing.T) {
//	_ = Init(&testConfig)
//	repo, err := New(RepoType("git"), url, "itools")
//	if err != nil {
//		t.Error(err)
//	}
//	err = repo.CheckoutToCommit("main", "87846d8941c92e4f2d5ceece0fc9465afb3c08fc")
//	if err != nil {
//		t.Error(err)
//	}
//	fmt.Println("CheckoutToCommit success")
//}

//func TestCheckoutToTag(t *testing.T) {
//	_ = Init(&testConfig)
//	repo, err := New(RepoType("git"), url, "itools")
//	if err != nil {
//		t.Error(err)
//	}
//	err = repo.CheckoutToTag("v1.0.0")
//	if err != nil {
//		t.Error(err)
//	}
//	fmt.Println("CheckoutToCommit success")
//}
