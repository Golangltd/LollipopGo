module LollipopGo

go 1.14

require (
	code.google.com/p/go.net/websocket v0.0.0-20200610144333-40b6804bddb4 // indirect
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.2.0
	github.com/gomodule/redigo v1.8.2
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/jinzhu/gorm v1.9.12
	github.com/name5566/leaf v0.0.0-20200516012428-8592b1abbbbe
	github.com/pkg/errors v0.9.1
	go.uber.org/zap v1.15.0
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859
	golang.org/x/sync v0.0.0-20200317015054-43a5402ce75a // indirect
)

replace (
	code.google.com/p/go.net/websocket => github.com/Golangltd/websocket_old v0.0.0-20200610144333-40b6804bddb4
	golang.org/x/net/websocket => github.com/Golangltd/websocket_old v0.0.0-20200610144333-40b6804bddb4
)
