package bootstrap

import (
	model "github.com/alainQtec/FileServer/v3/models"
	"github.com/alainQtec/FileServer/v3/models/scripts"
	"github.com/alainQtec/FileServer/v3/pkg/aria2"
	"github.com/alainQtec/FileServer/v3/pkg/auth"
	"github.com/alainQtec/FileServer/v3/pkg/cache"
	"github.com/alainQtec/FileServer/v3/pkg/cluster"
	"github.com/alainQtec/FileServer/v3/pkg/conf"
	"github.com/alainQtec/FileServer/v3/pkg/crontab"
	"github.com/alainQtec/FileServer/v3/pkg/email"
	"github.com/alainQtec/FileServer/v3/pkg/mq"
	"github.com/alainQtec/FileServer/v3/pkg/task"
	"github.com/alainQtec/FileServer/v3/pkg/wopi"
	"github.com/gin-gonic/gin"
	"io/fs"
	"path/filepath"
)

// Init 初始化启动
func Init(path string, statics fs.FS) {
	InitApplication()
	conf.Init(path)
	// Debug 关闭时 , 切换为生产模式
	if !conf.SystemConfig.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	dependencies := []struct {
		mode    string
		factory func()
	}{
		{
			"both",
			func() {
				scripts.Init()
			},
		},
		{
			"both",
			func() {
				cache.Init()
			},
		},
		{
			"slave",
			func() {
				model.InitSlaveDefaults()
			},
		},
		{
			"slave",
			func() {
				cache.InitSlaveOverwrites()
			},
		},
		{
			"master",
			func() {
				model.Init()
			},
		},
		{
			"both",
			func() {
				cache.Restore(filepath.Join(model.GetSettingByName("temp_path"), cache.DefaultCacheFile))
			},
		},
		{
			"both",
			func() {
				task.Init()
			},
		},
		{
			"master",
			func() {
				cluster.Init()
			},
		},
		{
			"master",
			func() {
				aria2.Init(false, cluster.Default, mq.GlobalMQ)
			},
		},
		{
			"master",
			func() {
				email.Init()
			},
		},
		{
			"master",
			func() {
				crontab.Init()
			},
		},
		{
			"master",
			func() {
				InitStatic(statics)
			},
		},
		{
			"slave",
			func() {
				cluster.InitController()
			},
		},
		{
			"both",
			func() {
				auth.Init()
			},
		},
		{
			"master",
			func() {
				wopi.Init()
			},
		},
	}

	for _, dependency := range dependencies {
		if dependency.mode == conf.SystemConfig.Mode || dependency.mode == "both" {
			dependency.factory()
		}
	}

}
