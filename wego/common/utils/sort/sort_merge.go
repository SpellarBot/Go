package sort

func part(L []int) ([]int, []int) {
	L1 := make([]int, 0)
	L2 := make([]int, 0)
	N := len(L)
	if N > 1 {
		n := N / 2
		L1 = L[0:n]
		L2 = L[n:N]
	}
	return L1, L2
}

func merge(L1 []int, L2 []int) []int {
	N1 := len(L1)
	N2 := len(L2)
	i, j := 0, 0
	var x, y int
	L := make([]int, 0)
	for {
		if (i >= N1) && (j >= N2) {
			break
		} else {
			if i < N1 && j < N2 {
				x = L1[i]
				y = L2[j]
				if x < y {
					L = append(L, x)
					i++
				} else {
					L = append(L, y)
					j++
				}
			} else {
				if i < N1 {
					L = append(L, L1[i])
					i++
				} else {
					L = append(L, L2[j])
					j++
				}
			}
		}
	}
	return L

}

func SortMerge(L []int) {
	N := len(L)
	if N > 1 {
		L1, L2 := part(L)
		SortMerge(L1)
		SortMerge(L2)
		L3 := merge(L1, L2)
		for i := 0; i < N; i++ {
			L[i] = L3[i]
		}
		return
	} else {
		return
	}
}
