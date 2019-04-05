package main

import (
	"fmt"
	"github.com/spaolacci/murmur3"
	"math/rand"
	"strconv"
)

func hash(id string) uint64 {
	h64Byte := murmur3.New64()
	h64Byte.Write([]byte(id))
	hash := h64Byte.Sum64()
	return hash
}

func main() {
	for i := 0; i < 10; i ++ {
		id := strconv.FormatUint(rand.Uint64(), 36) + strconv.FormatInt(int64(i), 8)
		buck0 := hash(id) % 1000000
		buck1 := hash(fmt.Sprintf("%d",buck0) + "1") % 1000000
		buck2 := hash(fmt.Sprintf("%d",buck0) + "2") % 1000000
		buck3 := hash(fmt.Sprintf("%d",buck0) + "3") % 1000000
		fmt.Println(buck0, buck1, buck2, buck3)
	}
}
