/*
Golang语言社区(www.Golang.Ltd)
作者：cserli
时间：2018年3月3日
*/

package controller

import (
	"sync"
)

// 控制器的结构体
type Tcontroller struct {
	MapController map[string]string
	mc            sync.RWMutex
}

// 申请对象
func NewController(paras ...interface{}) (m *Tcontroller) {

	return
}

// 注册路由
func (this *Tcontroller) RegisterTcontroller() {

	return
}

// 删除路由
func (this *Tcontroller) DeleteTcontroller() {

	return
}

// 修改路由
func (this *Tcontroller) ModifyTcontroller() {

	return
}
