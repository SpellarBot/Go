package MySort

import (
	"fmt"
)

func sift(a []int, k int, m int) {
	i := k
	j := 2*i + 1
	var x int
	for {
		if j > m {
			return
		}
		if j < m && a[j] < a[j+1] {
			j++
		}
		if a[i] > a[j] {
			return
		} else {
			x = a[i]
			a[i] = a[j]
			a[j] = x
			i = j
			j = 2*i + 1
		}
		//fmt.Println("----")
		//print_array(a)
	}

}

func SortHeap(a []int) {
	N := len(a)
	var x int
	for i := N / 2; i >= 1; i-- {
		fmt.Println("####")
		sift(a, i-1, N-1)
	}
	//print_array(a)
	for i := 1; i < N; i++ {
		x = a[0]
		a[0] = a[N-i]
		a[N-i] = x
		sift(a, 0, N-i-1)
		fmt.Println(i, "----")
		//print_array(a)
	}
	return
}
