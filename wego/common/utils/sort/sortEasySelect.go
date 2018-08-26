package sort

func SortEasySelect(a []interface{}, desc bool, compare func(a interface{}, b interface{}) bool) {
	N := len(a)
	var k int
	var x interface{}
	var m interface{}
	if desc {
		for i := 0; i < N; i++ {
			k = i
			m = a[i]
			for j := (i + 1); j < N; j++ {
				if !compare(m, a[j]) {
					k = j
					m = a[j]
				}
			}
			x = a[i]
			a[i] = a[k]
			a[k] = x
		}
	} else {
		for i := 0; i < N; i++ {
			k = i
			m = a[i]
			for j := (i + 1); j < N; j++ {
				if compare(m, a[j]) {
					k = j
					m = a[j]
				}
			}
			x = a[i]
			a[i] = a[k]
			a[k] = x
		}
	}

}
