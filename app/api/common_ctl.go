package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/wuzfei/go-helper/unit"
	"go-walle/app/internal/response"
	"go-walle/app/version"
	"runtime"
	"time"
)

type ServerInfo struct {
	Os   Os   `json:"os"`
	Cpu  Cpu  `json:"cpu"`
	Ram  Ram  `json:"ram"`
	Disk Disk `json:"disk"`
}

type Os struct {
	GOOS         string `json:"goos"`
	GOARCH       string `json:"goarch"`
	CPUNum       int    `json:"cpu_num"`
	Compiler     string `json:"compiler"`
	GoVersion    string `json:"go_version"`
	GoroutineNum int    `json:"goroutine_num"`
}

type Cpu struct {
	Cpus  []float64 `json:"cpus"`
	Cores int       `json:"cores"`
}

type Ram struct {
	Used        string  `json:"used"`
	Total       string  `json:"total"`
	UsedPercent float64 `json:"used_percent"`
}

type Disk struct {
	Used        string  `json:"used"`
	Total       string  `json:"total"`
	UsedPercent float64 `json:"used_percent"`
}

func getServerInfo() (_ *ServerInfo, err error) {
	s := ServerInfo{}
	s.Os = initOS()
	if s.Cpu, err = initCPU(); err != nil {
		return
	}
	if s.Ram, err = InitRAM(); err != nil {
		return
	}
	if s.Disk, err = initDisk(); err != nil {
		return
	}
	return &s, nil
}

func initOS() (o Os) {
	o.GOOS = runtime.GOOS
	o.GOARCH = runtime.GOARCH
	o.CPUNum = runtime.NumCPU()
	o.Compiler = runtime.Compiler
	o.GoVersion = runtime.Version()
	o.GoroutineNum = runtime.NumGoroutine()
	return o
}

func initCPU() (c Cpu, err error) {
	if cores, err := cpu.Counts(false); err != nil {
		return c, err
	} else {
		c.Cores = cores
	}
	if cpus, err := cpu.Percent(time.Duration(200)*time.Millisecond, true); err != nil {
		return c, err
	} else {
		c.Cpus = cpus
	}
	return c, nil
}

func InitRAM() (r Ram, err error) {
	if u, err := mem.VirtualMemory(); err != nil {
		return r, err
	} else {
		r.Used = unit.ByteFormat(int64(u.Used), 2)
		r.Total = unit.ByteFormat(int64(u.Total), 2)
		r.UsedPercent = u.UsedPercent
	}
	return r, nil
}

func initDisk() (d Disk, err error) {
	if u, err := disk.Usage("/"); err != nil {
		return d, err
	} else {
		d.Used = unit.ByteFormat(int64(u.Used), 2)
		d.Total = unit.ByteFormat(int64(u.Total), 2)
		d.UsedPercent = u.UsedPercent
	}
	return d, nil
}

type CommonCtl struct {
}

func (*CommonCtl) Version(ctx *gin.Context) {
	response.Success(ctx, version.Build)
}

func (*CommonCtl) ServiceInfo(ctx *gin.Context) {
	info, err := getServerInfo()
	response.Response(ctx, err, info)
}
