package sort

//冒泡排序
func SortPop(a []int) {
	N := len(a)
	var b int
	for i := 0; i < N; i++ {
		for j := i + 1; j < N; j++ {
			if a[i] > a[j] {
				b = a[i]
				a[i] = a[j]
				a[j] = b
			}
		}
	}
}
