package LollipopGo

import (
	"LollipopGo/LollipopGo/conf" // 配置文件
	"LollipopGo/LollipopGo/log"
)

func Run() {
	// logger
	if conf.LogLevel != "" {
		logger, err := log.New(conf.LogLevel, conf.LogPath, conf.LogFlag)
		if err != nil {
			panic(err)
		}
		log.Export(logger)
		defer logger.Close()
	}

	log.Release("Golang语言社区  LollipopGo %v starting up", version)
}
