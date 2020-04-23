package main

import (
	"github.com/nebulaim/telegramd/baselib/app"
	"github.com/nebulaim/telegramd/baselib/net2/examples/multi_proxy/server"
)

func main() {
	instance := &server.MultiProtoInsance{}
	// app.AppInstance(instance)
	app.DoMainAppInstance(instance)
}
