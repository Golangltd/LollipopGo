package SM

import "log"

// 有限状态机
type FSM struct {
	CurrentState int    // 当前状态数组的索引数值
	State []int         // 自定义状态
	// CallBcaks Callbacks // 回调函数
}

type Callback func()

type Callbacks map[string]Callback

func NewFSM(data []int) *FSM {
	if len(data) == 0{
		log.Println("create new FSM is fail")
		return nil
	}
	return &FSM{
		CurrentState: 0,
		State:data,
	}
}

func (this *FSM)NextTurn()  {
      if this!=nil{
      	if this.CurrentState <len(this.State)+1{
			this.CurrentState++
		}else {
			this.CurrentState = 0
		}
	  }else {
	  	log.Println("FSM is nil")
	  }
}

func (this *FSM)GetFSMState() int {
	 return this.State[this.CurrentState]
}

func (this *FSM)InitFSM()  {
	this = NewFSM(this.State)
}