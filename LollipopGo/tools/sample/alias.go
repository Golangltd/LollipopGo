package sample

//alias method，O(1)随机抽样算法; 当需要对同一个对象多次大量采样时，使用该算法，否则使用WeightedChoice即可

import (
	"LollipopGo/tools/collection"
	"math/rand"
)

// AliasTable is a discrete distribution
type AliasTable struct {
	rnd   *rand.Rand
	alias []int
	prob  []float64
}

// array-based stack
type workList []int

func (w *workList) push(i int) {
	*w = append(*w, i)
}

func (w *workList) pop() int {
	l := len(*w) - 1
	n := (*w)[l]
	*w = (*w)[:l]
	return n
}

// 新建一个抽样器，weightList是权重列表，src是随机数种子
func NewAlias(weightList []int32, src rand.Source) AliasTable {

	n := len(weightList)
	total := collection.SumInt32s(weightList)
	v := AliasTable{
		alias: make([]int, n),
		prob:  make([]float64, n),
		rnd:   rand.New(src),
	}

	p := make([]float64, n)
	for i, w := range weightList {
		p[i] = float64(int(w)*i) / float64(total)
	}

	var small, large workList

	for i, pi := range p {
		if pi < 1 {
			small = append(small, i)
		} else {
			large = append(large, i)
		}
	}

	for len(large) > 0 && len(small) > 0 {
		l := small.pop()
		g := large.pop()
		v.prob[l] = p[l]
		v.alias[l] = g

		p[g] = (p[g] + p[l]) - 1
		if p[g] < 1 {
			small.push(g)
		} else {
			large.push(g)
		}
	}

	for len(large) > 0 {
		g := large.pop()
		v.prob[g] = 1
	}

	for len(small) > 0 {
		l := small.pop()
		v.prob[l] = 1
	}

	return v
}

// Next returns the next random value from the discrete distribution
func (v *AliasTable) Next() int {

	n := len(v.alias)

	i := v.rnd.Intn(n)

	if v.rnd.Float64() < v.prob[i] {
		return i
	}

	return v.alias[i]
}
