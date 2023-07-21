package deploys

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wuzfei/go-helper/compress"
	"github.com/wuzfei/go-helper/path"
	"github.com/wuzfei/go-helper/slices"
	"github.com/zeebo/errs"
	"go-walle/app/global"
	"go-walle/app/model"
	"go-walle/app/pkg/repo"
	"go-walle/app/pkg/ssh"
	"go.uber.org/zap"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var ErrDeploy = errs.Class("deploy")

var ErrStopDeploy = ErrDeploy.New("终止发布任务")

var currentUser *user.User
var currentHost string

func init() {
	var err error
	currentUser, err = user.Current()
	if err != nil {
		panic("获取当前用户错误！")
	}
	currentHost, _ = os.Hostname()
}

type deployDirs struct {
	localWarehouseDir, localCodePackage, remoteReleaseDir, remoteReleasePackage, remoteRootLink string
}

type Task struct {
	model  *model.Task
	userId int64

	mux       *sync.Mutex
	ctx       context.Context
	cancel    context.CancelFunc
	doneError chan error
	isStop    chan struct{}
	started   bool
	stopped   bool

	env           []string
	deployDirs    *deployDirs
	deployPath    string
	deployPackage string
	targetRoot    string
	targetRelease string
}

func NewTask(model *model.Task, userId int64) *Task {
	return &Task{
		model:     model,
		userId:    userId,
		mux:       &sync.Mutex{},
		doneError: make(chan error),
	}
}

// check 检查基本状态是否可以发布上线
func (t *Task) check() error {
	if t.model.Status != model.TaskStatusAudit {
		return errors.New("任务未处于审核通过状态，无法发布")
	}
	if !t.model.Environment.Status.IsEnable() {
		return fmt.Errorf("该环境[%s]已不在允许发版本，请联系相关负责人处理", t.model.Environment.Name)
	}
	if !t.model.Project.Status.IsEnable() {
		return fmt.Errorf("该项目[%s]已不在允许发版本，请联系相关负责人处理", t.model.Project.Name)
	}
	if len(t.model.Servers) == 0 {
		return fmt.Errorf("该任务[%s]发布服务器为空，请联系相关负责人处理", t.model.Name)
	}
	return nil
}

func (t *Task) Start() (err error) {
	defer func() {
		if err != nil {
			removeDeployTask(t.model.ID)
			err = ErrDeploy.Wrap(err)
		}
	}()
	//检查基本状态
	err = t.check()
	if err != nil {
		return
	}
	//开始状态
	t.mux.Lock()
	if t.started {
		t.mux.Unlock()
		return nil
	}
	t.started = true
	t.mux.Unlock()

	//更新发布状态和版本
	t.model.Status = model.TaskStatusRelease
	t.model.Version = t.createReleaseVersion()
	err = global.DB.Select("status", "link_id").UpdateColumns(t.model).Error
	if err != nil {
		return
	}

	//启动发布协程
	t.ctx, t.cancel = context.WithCancel(context.Background())
	t.isStop = make(chan struct{}, 1)
	go func() {
		t.start()
	}()
	go func() {
		t.stop()
	}()
	return
}

func (t *Task) start() {
	var err error
loopFor:
	for _, f := range []func() error{t.prevDeploy, t.deploy, t.postDeploy, t.remoteRelease} {
		select {
		case <-t.ctx.Done():
			err = ErrStopDeploy
			break loopFor
		default:
			err = f()
			if err != nil {
				break loopFor
			}
		}
	}
	t.doneError <- err
}

func (t *Task) Stop() error {
	t.mux.Lock()
	started := t.started
	stopped := t.stopped
	t.mux.Unlock()
	if started && !stopped {
		t.cancel()
	}
	return nil
}

func (t *Task) IsStop() <-chan struct{} {
	t.mux.Lock()
	defer t.mux.Unlock()
	return t.isStop
}

func (t *Task) stop() {
	defer removeDeployTask(t.model.ID)
	doneErr := <-t.doneError
	close(t.doneError)
	t.mux.Lock()
	t.stopped = true
	close(t.isStop)
	t.mux.Unlock()
	t.model.Status = model.TaskStatusFinish
	if doneErr != nil {
		t.model.LastError = doneErr.Error()
		t.model.Status = model.TaskStatusReleaseFail
		if re, ok := doneErr.(RemoteErrs); ok {
			if re.HasSuccess() {
				t.model.Status = model.TaskStatusFinish
			}
		}
	}
	mb, _ := json.Marshal(t.model)

	if t.deployDirs.localCodePackage != "" {
		_ = os.RemoveAll(t.deployDirs.localCodePackage)
		_ = os.RemoveAll(t.deployDirs.localWarehouseDir)
	}

	if err := global.DB.Model(t.model).Select("status", "last_error").UpdateColumns(t.model).Error; err != nil {
		global.Log.Error("部署完成，更新数据库时出错", zap.ByteString("task_model", mb), zap.Error(doneErr), zap.Error(err))
	} else {
		global.Log.Debug("部署完成", zap.ByteString("task_model", mb))
	}
}

func (t *Task) remoteRelease() error {
	remoteErrs := make(RemoteErrs)
	wg := sync.WaitGroup{}
	for _, s := range t.model.Servers {
		wg.Add(1)
		go func(server *model.Server) {
			remoteErrs[server.ID] = t.remoteRun(server)
			wg.Done()
		}(s)
	}
	wg.Wait()
	return remoteErrs
}

// remoteRun 远程服务器执行部署
func (t *Task) remoteRun(server *model.Server) error {
	for _, f := range []func(server *model.Server) error{t.prevRelease, t.release, t.postRelease} {
		select {
		case <-t.ctx.Done():
			return ErrStopDeploy
		default:
			if err := f(server); err != nil {
				return err
			}
		}
	}
	return nil
}

// prevDeploy step1.检出代码前置操作
func (t *Task) prevDeploy() error {
	//1、检查仓库，
	_repo, err := t.getRepo()
	if err != nil {
		return errors.New("获取代码仓库错误：" + err.Error())
	}
	localDeployDir := filepath.Dir(_repo.Path())
	packageName := t.model.Version + ".tar.gz"
	t.deployDirs = &deployDirs{
		localWarehouseDir:    filepath.Join(localDeployDir, t.model.Version),
		localCodePackage:     filepath.Join(localDeployDir, packageName),
		remoteReleaseDir:     filepath.Join(t.model.Project.TargetReleases, t.model.Version),
		remoteReleasePackage: filepath.Join(t.model.Project.TargetReleases, packageName),
		remoteRootLink:       t.model.Project.TargetRoot,
	}
	//2、执行用户打包前命令
	commands := parseCommands(t.model.Project.PrevDeploy)
	for _, cmd := range commands {
		r := NewRecord(model.RecordTypePrevDeploy, t.model.ID, t.userId, cmd, nil, t.envs())
		if err := r.Run(t.ctx); err != nil {
			return err
		}
	}
	return nil
}

// deploy step2.检出代码
func (t *Task) deploy() error {
	//1、检出代码
	_repo, err := t.getRepo()
	if err != nil {
		return errors.New("获取代码仓库错误：" + err.Error())
	}
	if t.model.Tag != "" {
		err = _repo.CheckoutToTag(t.model.Tag)
	} else if t.model.Branch != "" && t.model.CommitId != "" {
		err = _repo.CheckoutToCommit(t.model.Branch, t.model.CommitId)
	} else {
		err = errors.New("发布分支选取错误")
	}
	if err != nil {
		return err
	}
	//2、复制发布版本代码到新目录，以便下面执行编译等操作
	if _, err = path.CopyDirToDir(t.deployDirs.localWarehouseDir, _repo.Path()); err != nil {
		return errors.New("检出代码失败：" + err.Error())
	}
	return err
}

// postDeploy step3.推送到服务器前的操作，比如下载依赖，编译等
func (t *Task) postDeploy() error {
	//1、在检出代码执行用户命令
	commands := parseCommands(t.model.Project.PostDeploy)
	for _, cmd := range commands {
		cmd = fmt.Sprintf("cd %s && %s", t.deployDirs.localWarehouseDir, cmd)
		r := NewRecord(model.RecordTypePostDeploy, t.model.ID, t.userId, cmd, nil, t.envs())
		if err := r.Run(t.ctx); err != nil {
			return err
		}
	}
	//2、打包代码
	st := time.Now()
	cmd := fmt.Sprintf("tar -zcvf %s -C %s", t.deployDirs.localCodePackage, t.deployDirs.localWarehouseDir)
	record := NewRecord(model.RecordTypePostDeploy, t.model.ID, t.userId, cmd, nil, nil)
	err := compress.PackMatch(t.deployDirs.localCodePackage, t.deployDirs.localWarehouseDir, t.getFileMatch())
	if err != nil {
		_err := "打包代码出错:" + err.Error()
		_ = record.Save(255, &_err, time.Since(st).Milliseconds())
		return err
	}
	_err := "success"
	_ = record.Save(0, &_err, time.Since(st).Milliseconds())
	return err
}

// prevRelease step4.推送代码到服务器前的操作
func (t *Task) prevRelease(server *model.Server) error {
	//解压程序包
	//_tarCmd := fmt.Sprintf("mkdir -p %s ", filepath.Dir(t.deployDirs.remoteReleasePackage))
	//r := NewRecord(model.RecordTypePrevRelease, t.model.ID, t.userId, _tarCmd, server, t.envs())
	//if err := r.Run(t.ctx); err != nil {
	//	return err
	//}
	//1、上传程序包
	st := time.Now()
	_saveCmd := fmt.Sprintf("scp -P%d %s@%s:%s %s:%s", server.Port, currentUser.Username, currentHost, t.deployDirs.localCodePackage, server.Hostname(), t.deployDirs.remoteReleasePackage)
	record := NewRecord(model.RecordTypePrevRelease, t.model.ID, t.userId, _saveCmd, server, nil)
	sftp, err := global.Ssh.NewSftp(ssh.ServerConfig{Host: server.Host, User: server.User, Port: server.Port})
	if err == nil {
		err = sftp.Copy(t.deployDirs.localCodePackage, t.deployDirs.remoteReleasePackage)
	}
	if err != nil {
		_err := "上传程序出错:" + err.Error()
		_ = record.Save(255, &_err, time.Since(st).Milliseconds())
		return err
	}
	_err := "success"
	_ = record.Save(0, &_err, time.Since(st).Milliseconds())

	//2、解压程序包
	_tarCmd := fmt.Sprintf("mkdir -p %s && tar -zxvf %s -C %s", t.deployDirs.remoteReleaseDir, t.deployDirs.remoteReleasePackage, t.deployDirs.remoteReleaseDir)
	r := NewRecord(model.RecordTypePrevRelease, t.model.ID, t.userId, _tarCmd, server, t.envs())
	if err := r.Run(t.ctx); err != nil {
		return err
	}
	//3、执行用户命令
	commands := parseCommands(t.model.Project.PrevRelease)
	for _, cmd := range commands {
		cmd = fmt.Sprintf("cd %s && %s", t.deployDirs.remoteReleaseDir, cmd)
		r := NewRecord(model.RecordTypePrevRelease, t.model.ID, t.userId, cmd, server, t.envs())
		if err := r.Run(t.ctx); err != nil {
			return err
		}
	}
	return nil
}

// release step5.部署程序
func (t *Task) release(server *model.Server) error {
	//1、获取上一个部署版本，保存下来
	cmd := fmt.Sprintf("[ -L %s ] && readlink %s || echo \"\"", t.deployDirs.remoteRootLink, t.deployDirs.remoteRootLink)
	record := NewRecord(model.RecordTypePrevRelease, t.model.ID, t.userId, cmd, server, t.envs())
	if err := record.Run(t.ctx); err != nil {
		return err
	}
	t.model.PrevVersion = record.Output()

	//2、部署代码，创建并替换源软连接
	tmpLink := fmt.Sprintf("%s_tmp", t.deployDirs.remoteRootLink)
	cmd = fmt.Sprintf("mkdir -p %s && ln -sfn %s %s", filepath.Dir(t.deployDirs.remoteRootLink), t.deployDirs.remoteReleaseDir, tmpLink)
	record = NewRecord(model.RecordTypePrevRelease, t.model.ID, t.userId, cmd, server, t.envs())
	if err := record.Run(t.ctx); err != nil {
		return err
	}

	cmd = fmt.Sprintf("mv -fT %s %s", tmpLink, t.deployDirs.remoteRootLink)
	record = NewRecord(model.RecordTypePrevRelease, t.model.ID, t.userId, cmd, server, t.envs())
	if err := record.Run(t.ctx); err != nil {
		return err
	}
	global.DB.Select("prev_version").UpdateColumns(t.model)
	return nil
}

// postRelease 6、执行部署完成功后用户相关命令
func (t *Task) postRelease(server *model.Server) error {
	commands := parseCommands(t.model.Project.PostRelease)
	for _, cmd := range commands {
		cmd = fmt.Sprintf("cd %s && %s", t.deployDirs.remoteRootLink, cmd)
		r := NewRecord(model.RecordTypePostRelease, t.model.ID, t.userId, cmd, server, t.envs())
		if err := r.Run(t.ctx); err != nil {
			return err
		}
	}
	return nil
}

func (t *Task) updateModel(task *model.Task) error {
	where := model.Task{ID: task.ID, UpdatedAt: t.model.UpdatedAt}
	return global.DB.Where(where).UpdateColumns(task).Error
}

func (t *Task) envs() *ssh.Envs {
	_envs := ssh.NewEnvsBySliceKV(parseCommands(t.model.Project.TaskVars))
	//_envs := NewEnvs()
	_envs.Add("PROJECT_ID", t.model.Project.ID)
	_envs.Add("PROJECT_NAME", t.model.Project.Name)
	_envs.Add("TASK_ID", t.model.ID)
	_envs.Add("TASK_NAME", t.model.Name)
	_envs.Add("DEPLOY_PATH", t.deployPath)
	_envs.Add("RELEASE_PATH", &t.model.Project.TargetRoot)
	return _envs
}

func (t *Task) createReleaseVersion() string {
	return fmt.Sprintf("%d_%d_%s", t.model.Project.ID, t.model.ID, time.Now().Format("20060102_150405"))
}

func (t *Task) getRepo() (repo.Repo, error) {
	return global.Repo.New(repo.TypeRepo(t.model.Project.RepoType), t.model.Project.RepoUrl, fmt.Sprintf("%d", t.model.Project.ID))
}

func (t *Task) getFileMatch() compress.Match {
	regs := strings.Split(t.model.Project.Excludes, "\n")
	if len(regs) > 0 {
		regs = slices.Map(regs, func(item string, k int) string {
			return strings.TrimSpace(item)
		})
		if t.model.Project.IsInclude == 1 {
			return compress.FileMatch(regs...)
		}
		return compress.ReFileMatch(regs...)
	}
	return nil
}

// parseCommands 解析命令，支持'#'，'//'的行注释
func parseCommands(commands string) []string {
	res := make([]string, 0)
	commands = strings.TrimSpace(commands)
	if commands == "" {
		return res
	}
	arr := strings.Split(commands, "\n")
	for _, v := range arr {
		v = strings.TrimSpace(v)
		if v == "" || v[:1] == "#" || (len(v) > 1 && v[:2] == "//") {
			continue
		}
		res = append(res, v)
	}
	return res
}
