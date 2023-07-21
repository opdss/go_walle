package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wuzfei/cfgstruct/cfgstruct"
	"github.com/wuzfei/cfgstruct/process"
	"github.com/wuzfei/go-helper/path"
	"go-walle/app/api"
	"go-walle/app/global"
	"go-walle/app/migration"
	db2 "go-walle/app/pkg/db"
	log3 "go-walle/app/pkg/log"
	"go-walle/app/version"
	log2 "log"
	"os"
	"path/filepath"
)

//go:generate stringer -type ErrCode -linecomment ./app/internal/errcode

//go:embed web/dist/*
var web embed.FS

// 多加这个是因为前端打包的资源里面包含了_开头的文件
//
//go:embed web/dist/assets/*
var webAssets embed.FS

var (
	runCfg   global.Config
	setupCfg global.Config
)

var (
	configFile string
	rootCmd    = &cobra.Command{
		Use:   "gowalle",
		Short: "简单部署系统",
	}
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "运行",
		RunE:  cmdRun,
	}
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "查看当前所有配置",
		RunE:  cmdConfig,
	}
	migrationCmd = &cobra.Command{
		Use:   "migration",
		Short: "数据库迁移命令",
		RunE:  cmdMigration,
	}
	setupCmd = &cobra.Command{
		Use:         "setup",
		Short:       "初始化配置",
		RunE:        cmdSetup,
		Annotations: map[string]string{"type": "setup"},
	}
)

func main() {
	log2.Println(version.Build.String())
	defaultConfig := path.ApplicationDir("go-walle", process.DefaultCfgFilename)
	cfgstruct.SetupFlag(rootCmd, &configFile, "config", defaultConfig, "配置文件")
	//根据环境读取默认配置
	defaults := cfgstruct.DefaultsFlag(rootCmd)
	//当前程序所在目录
	currentDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	rootDir := cfgstruct.ConfigVar("ROOT", currentDir)
	//设置系统的HOME变量
	envHome := cfgstruct.ConfigVar("HOME", os.Getenv("HOME"))
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(migrationCmd)
	rootCmd.AddCommand(configCmd)
	process.Bind(runCmd, &runCfg, defaults, cfgstruct.ConfigFile(configFile), envHome, rootDir)
	process.Bind(migrationCmd, &runCfg, defaults, cfgstruct.ConfigFile(configFile), envHome, rootDir)
	process.Bind(configCmd, &runCfg, defaults, cfgstruct.ConfigFile(configFile), envHome, rootDir)
	process.Bind(setupCmd, &setupCfg, defaults, cfgstruct.ConfigFile(configFile), envHome, cfgstruct.SetupMode(), rootDir)
	process.Exec(rootCmd)
}

// cmdRun 运行
func cmdRun(cmd *cobra.Command, args []string) (err error) {
	ctx, _ := process.Ctx(cmd)
	runCfg.Init()
	apiServer := api.NewServer(&runCfg, &web, &webAssets)
	return apiServer.Run(ctx)
}

// cmdSetup 初始化数据库
func cmdSetup(cmd *cobra.Command, args []string) error {
	return process.SaveConfig(cmd, configFile)
}

// cmdConfig 查看系统配置
func cmdConfig(cmd *cobra.Command, args []string) error {
	fmt.Printf("当前运行环境：[%s]\n", cfgstruct.DefaultsType())
	fmt.Println("当前配置文件路径：", configFile)
	output, _ := json.MarshalIndent(runCfg, "", " ")
	fmt.Println(string(output))
	return nil
}

// cmdMigration 数据库迁移初始化
func cmdMigration(cmd *cobra.Command, args []string) error {
	_log := log3.NewLog(&runCfg.Log)
	db, err := db2.NewGormDB(&runCfg.Db, _log)
	if err != nil {
		return err
	}
	dsn, _ := runCfg.Db.GetDsn()
	fmt.Println("运行数据库：", dsn)
	mr := migration.NewMigration(_log.Named("migration"), db)
	return mr.Setup()
}
