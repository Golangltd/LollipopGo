/*
Golang语言社区(www.Golang.Ltd)
作者：cserli
时间：2018年3月3日
*/

package API

//---单例函数

//---排序函数（支持多维排序）
//插入排序
//1从第一个元素开始，该元素可以认为已经被排序
//2取出下一个元素，在已经排序的元素序列中从后向前扫描
//3如果该元素（已排序）大于新元素，将该元素移到下一位置
//4重复步骤3，直到找到已排序的元素小于或者等于新元素的位置
//5将新元素插入到该位置后
//6重复步骤2~5
func insertSort(p []int) {

	for i := 1; i < len(p); i++ {

		for j := i - 1; j >= 0 && p[j+1] < p[j]; j-- {

			p[j+1], p[j] = p[j], p[j+1]

		}

	}

}

//冒泡排序
//算法描述
//1比较相邻的元素。如果第一个比第二个大，就交换他们两个。
//2对每一对相邻元素作同样的工作，从开始第一对到结尾的最后一对。这步做完后，最后的元素会是最大的数。
//3针对所有的元素重复以上的步骤，除了最后一个。
//4持续每次对越来越少的元素重复上面的步骤，直到没有任何一对数字需要比较
func bubbleSort(p []int) {

	for i := 1; i < len(p)-1; i++ {

		for j := 0; j < len(p)-i; j++ {

			if p[j] > p[j+1] {

				p[j+1], p[j] = p[j], p[j+1]

			}

		}

	}

}

// 使用
//	length := len(list)
//	for root := length/2 - 1; root >= 0; root-- {
//		sort(list, root, length)
//	} //第一次建立大顶堆
//	for i := length - 1; i >= 1; i-- {
//		list[0], list[i] = list[i], list[0]
//		sort(list, 0, i)
//	} //调整位置并建并从第一个root開始建堆.假设不明确为什么,大家多把图画几遍就应该明朗了
//	fmt.Println(list)

type Person struct {
	FirmsName  string // 商家的名字
	XCName     string // 现场的名字
	GameName   string // 游戏的名字
	OpenID     string // 谁获取了积分
	IJiFen     uint32 // 获得的积分
	Instertime string // 获取积分的时间
}

func sortbg(list []Person, root, length int) {
	for {
		child := 2*root + 1
		if child >= length {
			break
		}
		if child+1 < length && list[child].IJiFen > list[child+1].IJiFen {
			child++ //这里重点讲一下,就是调整堆的时候,以左右孩子为节点的堆可能也须要调整
		}
		if list[root].IJiFen < list[child].IJiFen {
			return
		}
		list[root], list[child] = list[child], list[root]
		root = child
	}
}
