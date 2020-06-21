package LollipopGo

import (
	"LollipopGo/log"
	"flag"
)

func Run() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "./log")
	flag.Set("v", "3")
	flag.Parse()
	log.Release("Golang语言社区  LollipopGo %v starting up", Version)
}
