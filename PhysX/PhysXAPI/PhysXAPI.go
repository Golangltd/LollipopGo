package PhysXAPI

import (
	"LollipopGo/log"
	"runtime"
	"sync"
	"syscall"
)

/*
   PhysX API 中文翻译地址：www.PhysXAPI.Com 网站开放时间预计：2021年1月1日
   发起人：彬哥
   微信号：cserli
*/

const (
	Linux   string = "linux"
	Windows string = "windows"
)

type LollipopGoPhysX struct {
	BSwitch    bool              // 是否开启 PhysX
	OsType     string            // 系统类型,linux和windows
	LoLock     sync.RWMutex      // 读写锁
	LoLazyDLL  *syscall.LazyDLL  // 系统调用的指针
	LoLazyProc *syscall.LazyProc // 系统调用函数指针
}

// 生成 NewLollipopGoPhysX的结构数据信息
func NewLollipopGoPhysX(bSwitch bool) *LollipopGoPhysX {
	return &LollipopGoPhysX{
		BSwitch: bSwitch,
		OsType:  runtime.GOOS,
	}
}

// 是否开启PhysX
func (this *LollipopGoPhysX) PhysXSwitch() {
	if !this.CheckPtr() {
		return
	}
	if this.BSwitch {
		this.LoadPhysXLibrary()
	}
}

// 加载数据信息
func (this *LollipopGoPhysX) LoadPhysXLibrary() {
	if !this.CheckPtr() {
		return
	}
	if this.OsType == Linux {
		// LINUX系统
	}
	if this.OsType == Windows {
		// windows系统
		dll32 := syscall.NewLazyDLL("./DllAndSo/LollipopGo.dll")
		log.Debug("call dll:", dll32.Name)
		this.LoLazyDLL = dll32
	}
}

// 获取相应的函数指针
// ret, _, _ :=g.Call(uintptr(4),uintptr(8))
func (this *LollipopGoPhysX) GetFnPtr(FnName string) {
	if !this.CheckPtr() {
		return
	}
	this.LoLazyProc = this.LoLazyDLL.NewProc(FnName)
}

// 检查是否错误
func (this *LollipopGoPhysX) CheckPtr() bool {
	if this == nil {
		log.Error("LollipopGo ptr is nil")
		return false
	}
	return true
}
