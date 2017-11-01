package abtest

import (
	"fmt"
	"testing"
	"strconv"
	"math/rand"
	"time"
)

func Test(t *testing.T) {
	file := "abtest_example.xml"

	abtestConf := ABtestConfig{}

	abtestConf.Init(file, 10)

	runNum := 100
	for runNum > 0 {
		runNum--

		mutex_defaultCnt := 0
		mutex_t1 := 0
		mutex_t2 := 0
		mutex_t3 := 0


		layer_defaultCnt := 0
		layer_t1 := 0
		layer_t2 := 0
		layer_t3 := 0

		cf_defaultCnt := 0
		cf_t1 := 0

		mutex_none := 0
		layer_none := 0
		cf_none := 0

		i := 200000
		for i > 0 {
			i--

			uid := strconv.FormatUint(rand.Uint64(), 2) +  strconv.FormatInt(int64(i), 10)
			//fmt.Println(uid)


			name1, ok1 := abtestConf.GetMutexABTag(uid)

			if ok1 {
				//fmt.Println(fmt.Sprintf("uid=%s\ttag=%s", uid, name1))
				switch name1 {
				case "default":
					mutex_defaultCnt++
				case "t1":
					mutex_t1++
				case "t2":
					mutex_t2++
				case "t3":
					mutex_t3++
				}
			} else {
				mutex_none++
			}

			name2, ok2 := abtestConf.GetABTag(uid, "ctr")
			if ok2 {
				//fmt.Println(fmt.Sprintf("uid=%s\ttag=%s", uid, name2))
				switch name2 {
				case "default":
					layer_defaultCnt++
				case "t1":
					layer_t1++
				case "t2":
					layer_t2++
				case "t3":
					layer_t3++
				}
			} else {
				layer_none++
			}

			name3, ok3 := abtestConf.GetABTag(uid, "cf")
			if ok3 {
				//fmt.Println(fmt.Sprintf("uid=%s\ttag=%s", uid, name3))
				switch name3 {
				case "default":
					cf_defaultCnt++
				case "t1":
					cf_t1++
				}
			} else {
				cf_none++
			}
		}

		fmt.Println("mutex ==================")
		fmt.Println("mutex_defaultCnt: ", mutex_defaultCnt)
		fmt.Println("mutex_t1: ", mutex_t1)
		fmt.Println("mutex_t2: ", mutex_t2)
		fmt.Println("mutex_t3: ", mutex_t3)

		fmt.Println("layer_ctr ==================")
		fmt.Println("layer_defaultCnt: ", layer_defaultCnt)
		fmt.Println("layer_t1: ", layer_t1)
		fmt.Println("layer_t2: ", layer_t2)
		fmt.Println("layer_t3: ", layer_t3)

		fmt.Println("layer_cf ==================")
		fmt.Println("cf_defaultCnt: ", cf_defaultCnt)
		fmt.Println("cf_t1: ", cf_t1)

		fmt.Println("layer_cf ==================")
		fmt.Println("mutex_none: ", mutex_none)
		fmt.Println("layer_none: ", layer_none)
		fmt.Println("cf_none: ", cf_none)

		time.Sleep(2 * time.Second)
	}

}