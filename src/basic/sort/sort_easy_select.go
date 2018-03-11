package sort

func SortEasySelect(a []int) {
	N := len(a)
	var k int
	var x int
	var m int
	for i := 0; i < N; i++ {
		k = i
		m = a[i]
		for j := (i + 1); j < N; j++ {
			if m > a[j] {
				k = j
				m = a[j]
			}
		}
		x = a[i]
		a[i] = a[k]
		a[k] = x
	}
}
