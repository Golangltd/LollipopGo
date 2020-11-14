package LollipopGo

import (
	"LollipopGo/log"
	"flag"
)

func init()  {
	log.Release("Golang语言社区  LollipopGo %v starting up", Version)
}

func Run() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "./log")
	flag.Set("v", "3")
	flag.Parse()
}
