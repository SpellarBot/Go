package abtest

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func BenchmarkAbtest(b *testing.B) {
	b.StopTimer()
	file := "abtest_test.xml"
	abtestConf := ABtestConfig{}
	abtestConf.Init(file, 10)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		uid := strconv.FormatUint(rand.Uint64(), 2) + strconv.FormatInt(int64(i), 10)
		abtestConf.GetTag(uid)
	}
}

func TestAbtest(t *testing.T) {
	file := "abtest_test.xml"
	abtestConf := ABtestConfig{}
	er := abtestConf.Init(file, 10)
	N := 1000000
	if er == nil {
		abtestConf.Print()
		mutexNum := make(map[string]int)
		layerNum := make(map[string]int)
		t1 := time.Now().UnixNano()
		for i := 0; i < N; i++ {
			uid := strconv.FormatUint(rand.Uint64(), 2) + strconv.FormatInt(int64(i), 10)
			tagType, Type, err := abtestConf.GetTag(uid)
			if err == true {
				if tagType == "mutex" {
					if num, ok := mutexNum[Type]; ok {
						mutexNum[Type]++
					} else {
						mutexNum[Type] = num
					}
				} else {
					if num, ok := layerNum[Type]; ok {
						layerNum[Type]++
					} else {
						layerNum[Type] = num
					}
				}
			}

		}
		t2 := time.Now().UnixNano()
		fmt.Println(fmt.Sprintf("Cost All Time: %5d s", (t2-t1)/1000000000))
		fmt.Println(fmt.Sprintf("Cost Eve Time: %5d ns", (t2-t1)/int64(N)))

		for key, value := range layerNum {
			fmt.Println(fmt.Sprintf("%30s %d", key, value))
		}

		for key, value := range mutexNum {
			fmt.Println(fmt.Sprintf("%30s %d", key, value))
		}
	}
}

// go test
// go test -v
// go test -v -run="TestAbtest" -count 5
// go test -bench="BenchmarkAbtest" -count 5
