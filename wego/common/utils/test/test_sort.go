package main

import "fmt"

func Sort(a []float64,videos[]string, desc bool) {
	N := len(a)
	var b float64
	var video string
	if desc == true{
		for i := 0; i < N; i++ {
			for j := i + 1; j < N; j++ {
				if a[i] < a[j] {
					b = a[i]
					a[i] = a[j]
					a[j] = b
					video =videos[i]
					videos[i] = videos[j]
					videos[j] = video
				}
			}
		}
	}else{
		for i := 0; i < N; i++ {
			for j := i + 1; j < N; j++ {
				if a[i] > a[j] {
					b = a[i]
					a[i] = a[j]
					a[j] = b
					video =videos[i]
					videos[i] = videos[j]
					videos[j] = video
				}
			}
		}
	}
}

func main(){
	videos := []string{"4.5","7.8","1.1","2.9","5.9","100","0","0","3.1","0.3"}
	a := []float64{4.5,7.8,1.1,2.9,5.9,100,0,0,3.1,0.3}
	Sort(a, videos, false)
	fmt.Println(a)
	fmt.Println(videos)
	C := make([]string,0)
	C =append(C,"dede")

}