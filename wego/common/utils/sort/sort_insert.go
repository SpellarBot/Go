package sort

//插入排序
func SortInsert(a []int) {
	for i, k := range a {
		if i == 0 {
			continue
		}
		//print_array(a[0:i])
		//f.Println(i+1, k)
		insert(a[0:i+1], i+1, k)
	}

}

// a[N-1] = x, a[0:N-2] in order
func insert(a []int, N int, x int) {
	var k int
	for k = 0; k < N-1; k++ {
		if x < a[k] {
			break
		}
	}
	//f.Println("The location is", k)
	for i := N - 1; i > k; i-- {
		//f.Println(a[i], a[i-1])
		a[i] = a[i-1]

	}
	a[k] = x
}
