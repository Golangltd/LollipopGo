package leaf

import (
	"os"
	"os/signal"

	"github.com/name5566/leaf/cluster"
	"github.com/name5566/leaf/conf"
	"github.com/name5566/leaf/console"
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/module"
)

// 注册和运行模块信息
func Run(mods ...module.Module) {
	// logger
	if conf.LogLevel != "" {
		logger, err := log.New(conf.LogLevel, conf.LogPath, conf.LogFlag)
		if err != nil {
			panic(err)
		}
		log.Export(logger)
		defer logger.Close()
	}
	log.Release("Golang语言社区  LeafLtd %v starting up", version)

	// module 模块
	for i := 0; i < len(mods); i++ {
		module.Register(mods[i])
	}
	module.Init()

	cluster.Init()

	console.Init()

	// close 关闭
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	log.Release("Leaf closing down (signal: %v)", sig)
	console.Destroy()
	cluster.Destroy()
	module.Destroy()
}
