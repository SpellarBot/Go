package main

import (
	"fmt"
	"math/rand"
	"time"
	"wego/common/utils/sort"
)

type data []int

func (d data) Len() int {
	return len(d)
}

func (d data) Compare(i, j int) bool {
	if d[i] > d[j] {
		return true
	} else {
		return false
	}
}
func (d data) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func main() {
	rand.Seed(time.Now().Unix())
	var d []int
	for i := 0; i < 10000; i++ {
		d = append(d, rand.Intn(1000))
	}
	a := data(d)
	// fmt.Println(d)
	t1 := time.Now().UnixNano()
	// sort.SortEasySelect(a, false)
	// sort.SortPop(a, false)
	// sort.SortHeap(a, true)
	//sort.SortInsert(a, true)
	sort.SortQuick(a, true)
	// sort.SortMerge(a, false)
	// gsort.Sort(gsort.IntSlice(d))
	t2 := time.Now().UnixNano()
	fmt.Println((t2 - t1) / (1000000))
	// fmt.Println(d)
	// a1 := data([]int{1, 4, 7, 2, 3, 8})
	// fmt.Println(a1)
	// sort.Merge(a1, 0, 2, 5)
	// fmt.Println(a1)
}
