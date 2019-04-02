package main //main函数
import "fmt" //相当于#include
func main() {
	ar := [10]int{9, 8, 6, 4, 2, 7, 1, 3, 0, 5}
	num := len(ar)             //:=自动匹配变量类型
	for i := 0; i < num; i++ { //花括号{必须在这一行 if也一样
		for j := i + 1; j < num; j++ {
			if ar[i] > ar[j] { // 排序的方式
				ar[i], ar[j] = ar[j], ar[i] //相比某语言不需要临时交换变量
			}
		}
	}
	fmt.Println(ar)
}
