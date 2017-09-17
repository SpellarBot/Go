package MySort

//快速排序

import (
	"fmt"
)

func partition(a []int, first int, end int) int {
	m := first
	n := end
	var x int
	for {
		fmt.Println(m, n)
		if m >= n {
			break
		}
		for {
			if m < n && a[m] <= a[n] {
				n--
			} else {
				break
			}
		}
		if m < n {
			x = a[m]
			a[m] = a[n]
			a[n] = x
			m++
		}

		for {
			if m < n && a[m] <= a[n] {
				m++
			} else {
				break
			}
		}
		if m < n {
			x = a[m]
			a[m] = a[n]
			a[n] = x
			n--
		}
	}
	return m
}

func sort_quick(a []int) {
	N := len(a)
	if N <= 1 {
		return
	}
	i := partition(a, 0, N-1)
	//fmt.Println(i)
	sort_quick(a[0 : i+1])
	sort_quick(a[i+1 : N])
}
