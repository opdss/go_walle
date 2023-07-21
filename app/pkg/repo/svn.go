package repo

type SvnConfig struct {
	Username string `help:"svn帐号" default:""`
	Password string `help:"svn密码" default:""`
}

type Svn struct {
	config *SvnConfig
	path   string
	url    string
}

func NewSvn(cfg *SvnConfig, url string, path string) (*Svn, error) {
	return &Svn{config: cfg, url: url, path: path}, nil
}

func (srv *Svn) Tags() ([]Tag, error) {
	return nil, ErrRepoSvn.New("todo")
}
func (srv *Svn) Commits(branch string) ([]Commit, error) {
	return nil, ErrRepoSvn.New("todo")
}
func (srv *Svn) Branches() ([]Branch, error) {
	return nil, ErrRepoSvn.New("todo")
}
func (srv *Svn) CheckoutToBranch(branch string) error {
	return ErrRepoSvn.New("todo")
}
func (srv *Svn) CheckoutToCommit(branch, commit string) error {
	return ErrRepoSvn.New("todo")
}
func (srv *Svn) CheckoutToTag(tag string) error {
	return ErrRepoSvn.New("todo")
}
func (srv *Svn) Path() string {
	return srv.path
}
func (srv *Svn) Type() TypeRepo {
	return SvnRepo
}
