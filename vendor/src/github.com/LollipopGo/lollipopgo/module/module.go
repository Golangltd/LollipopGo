package module

import (
	"runtime"
	"sync"

	"github.com/LollipopGo/lollipopgo/conf"
	"github.com/LollipopGo/lollipopgo/log"
)

// 外部接口
type Module interface {
	OnInit()
	OnDestroy()
	Run(closeSig chan bool)
}

// 内部结构
type module struct {
	mi       Module
	closeSig chan bool
	wg       sync.WaitGroup
}

// 定义一个模块的数组
var mods []*module

// 注册模块；也就是把数据保存起来（缓存起来）
func Register(mi Module) {
	m := new(module)
	m.mi = mi
	m.closeSig = make(chan bool, 1)

	mods = append(mods, m)
}

// 每个模块话的初始化
func Init() {
	for i := 0; i < len(mods); i++ {
		mods[i].mi.OnInit() // 每个模块实现自己的初始化函数
	}

	for i := 0; i < len(mods); i++ {
		m := mods[i]
		m.wg.Add(1)
		go run(m)
	}
}

// 关闭所有模块
func Destroy() {
	for i := len(mods) - 1; i >= 0; i-- {
		m := mods[i]
		m.closeSig <- true
		m.wg.Wait()
		destroy(m)
	}
}

func run(m *module) {
	m.mi.Run(m.closeSig)
	m.wg.Done()
}

func destroy(m *module) {
	defer func() {
		if r := recover(); r != nil {
			if conf.LenStackBuf > 0 {
				buf := make([]byte, conf.LenStackBuf)
				l := runtime.Stack(buf, false)
				log.Error("%v: %s", r, buf[:l])
			} else {
				log.Error("%v", r)
			}
		}
	}()

	m.mi.OnDestroy()
}
