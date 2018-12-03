package main

/*
Go语言调用C函数的方式：
1 嵌套到工程的代码当中，
2 将工程的中的C语言代码，分割出去，写成单独的文件。
3 将C语言写成DLL形式调用

<conio.h>
conio是Console Input/Output（控制台输入输出）的简写，
其中定义了通过控制台进行数据输入和数据输出的函数，
主要是一些用户通过按键盘产生的对应操作，比如getch()函数等等。
*/

/*
#include <windows.h>
#include <conio.h>

void gotoxy(int x ,int y)
{
	HANDLE handle = GetStdHandle(STD_OUTPUT_HANDLE);                        // 获取控制台句柄
    SetConsoleTextAttribute(handle, FOREGROUND_INTENSITY | FOREGROUND_RED); // 设置为红色
    COORD c;
    c.X=x,c.Y=y;
    SetConsoleCursorPosition(GetStdHandle(STD_OUTPUT_HANDLE),c);
}
// 从键盘获取一次按键，但不显示到控制台
int direct()
{
    return _getch();
}
// 控制台清屏数据
void system_cls()
{
	system("cls");
}
*/
import "C"

import (
	"flag"
	"fmt"
	_ "glog-master"
	"math/rand"
	"os"
	"time"

	"code.google.com/p/go.net/websocket"
)

var addr = flag.String("addr", "127.0.0.1:8893", "http service address")
var connbak *websocket.Conn
var bkaiguan chan bool

// 初始化操作
func init() {
	// 日志初始化
	flag.Set("alsologtostderr", "true") // 日志写入文件的同时，输出到stderr
	flag.Set("log_dir", "./log")        // 日志文件保存目录
	flag.Set("v", "3")                  // 配置V输出的等级。
	flag.Parse()
	bkaiguan = make(chan bool) // 开关
	// 初始化网络信息
	if initNet() {
		// 匹配 对战操作
		// initMatch(connbak)
		// initbak()  ----bak

		return
	}
	panic("链接服务器失败！！！")
	return
}

func initNet() bool {

	fmt.Println(" 用户客户端客户端模拟！")
	url := "ws://" + *addr + "/GolangLtdSnake"
	connbak, err := websocket.Dial(url, "", "test://golang/")
	if err != nil {
		fmt.Println("err:", err.Error())
		return false
	}
	// 协程支持  --接受线程操作
	go GameServerReceive(connbak)
	// 1 登录 协议
	initLogin(connbak)
	// 2 进入 游戏
	initMatch(connbak)
	return true
}

// 表示光标的位置
type loct struct {
	i, j int
}

// 数据定义
var (
	area = [20][20]byte{} // 记录了蛇、食物的信息
	food bool             // 当前是否有食物
	lead byte             // 当前蛇头移动方向
	head loct             // 当前蛇头位置
	tail loct             // 当前蛇尾位置
	size int              // 当前蛇身长度
)

// 随机生成一个位置，来放置食物import "C" // go中可以嵌入C语言的函数
func place() loct {
	k := rand.Int() % 400
	return loct{k / 20, k % 20}
}

// 用来更新控制台的显示，在指定位置写字符，使用错误输出避免缓冲
func draw(p loct, c byte) {
	C.gotoxy(C.int(p.i*2+4), C.int(p.j+2))
	fmt.Fprintf(os.Stderr, "%c", c)
}

func initbak() {

	// 初始化蛇的位置和方向、首尾；初始化随机数
	head, tail = loct{4, 4}, loct{4, 4}
	lead, size = 'R', 1
	area[4][4] = 'H'
	rand.Seed(int64(time.Now().Unix()))

	// 清屏操作
	C.system_cls()

	// 输出初始画面
	fmt.Fprintln(os.Stderr, `
#-----------------------------------------#
|                                         |
|                                         |
|                                         |
|                                         |
|         *                               |
|                                         |
|                                         |
|                                         |
|                                         |
|                                         |
|                                         |
|                                         |
|                                         |
|                                         |
|                                         |
|                                         |
|                                         |
|                                         |
|                                         |
|                                         |
#-----------------------------------------#
         字节教育  www.ByteEdu.com
`)

	// 我们使用一个单独的go程来捕捉键盘的动作，因为是单独的，不怕阻塞
	go func() {
		for { // 函数只写入lead，外部只读取lead，无需设锁
			switch byte(C.direct()) {
			case 72:
				lead = 'U'
			case 75:
				lead = 'L'
			case 77:
				lead = 'R'
			case 80:
				lead = 'D'
			case 32:
				lead = 'P'
			}
		}
	}()
}

func main() {

	// 主程序
	for {

		// 程序更新周期，400毫秒
		time.Sleep(time.Millisecond * 400)

		// 暂停，还是要有滴
		if lead == 'P' {
			continue
		}

		// 放置食物
		if !food {
			give := place()
			if area[give.i][give.j] == 0 { // 食物只能放在空闲位置
				area[give.i][give.j] = 'F'
				draw(give, '$') // 绘制食物
				food = true
			}
		}

		// 我们在蛇头位置记录它移动的方向
		area[head.i][head.j] = lead

		// 根据lead来移动蛇头
		switch lead {
		case 'U':
			head.j--
		case 'L':
			head.i--
		case 'R':
			head.i++
		case 'D':
			head.j++
		}

		// 判断蛇头是否出界
		if head.i < 0 || head.i >= 20 || head.j < 0 || head.j >= 20 {
			C.gotoxy(0, 23) // 让光标移动到画面下方
			// 玩家死亡  直接处理
			break // 跳出死循环
		}

		// 获取蛇头位置的原值，来判断是否撞车，或者吃到食物
		eat := area[head.i][head.j]
		if eat == 'F' { // 吃到食物
			food = false
			// 增加蛇的尺寸，并且不移动蛇尾
			size++
		} else if eat == 0 { // 普通移动

			draw(tail, ' ') // 擦除蛇尾

			// 注意我们记录了它移动的方向
			dir := area[tail.i][tail.j]

			// 我们需要擦除蛇尾的记录
			area[tail.i][tail.j] = 0

			// 移动蛇尾
			switch dir {
			case 'U':
				tail.j--
			case 'L':
				tail.i--
			case 'R':
				tail.i++
			case 'D':
				tail.j++
			}
		} else { // 撞车了
			C.gotoxy(0, 23)
			break
		}
		draw(head, '*') // 绘制蛇头
	}

	// 收尾了
	switch {
	case size < 22:
		fmt.Fprintf(os.Stderr, "Faild! You've eaten %d $\\n", size-1)
	case size < 42:
		fmt.Fprintf(os.Stderr, "Try your best! You've eaten %d $\\n", size-1)
	default:
		fmt.Fprintf(os.Stderr, "Congratulations! You've eaten %d $\\n", size-1)
	}
}
