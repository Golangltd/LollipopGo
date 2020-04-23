package http

import (
	"net/http"
	"time"
)

var DefaultClient *http.Client

func init() {
	clt := *http.DefaultClient
	clt.Timeout = time.Second * 5
	DefaultClient = &clt
}
