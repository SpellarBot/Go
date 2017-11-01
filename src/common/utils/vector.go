// 二相箔向量空间常用计算

package utils

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"encoding/binary"
)

const DimCnt int = 200 // 重要常量！向量空间的维度

func VecCompute(vecMap *map[string][DimCnt]float64, target []float64, ids []string, posBegin int, posEnd int) bool {
	if posBegin > len(ids) || posBegin < 0 ||
		posEnd > len(ids) || posEnd < 0 ||
		posBegin > posEnd {
		return false
	}
	ret := false
	for i := posBegin; i <= posEnd; i++ {
		vec, ok := (*vecMap)[ids[i]]
		if ok {
			for j := 0; j < DimCnt; j++ {
				target[j] += vec[j]
			}
			ret = true
		}
	}
	vecLen := 0.0
	for i := 0; i < DimCnt; i++ {
		vecLen += target[i] * target[i]
	}
	vecLen = math.Sqrt(vecLen)
	for i := 0; i < DimCnt; i++ {
		target[i] /= vecLen
	}
	return ret
}

func DislikeCompute(dislikeVec *[10][DimCnt]float64, v *[DimCnt]float64, scoreMin int, cnt int, found int) bool {
	ret := false
	if found > 0 {
		for i := 0; i < cnt; i++ {
			dist := 0.0
			for j := 0; j < DimCnt; j++ {
				dist += (*dislikeVec)[i][j] * (*v)[j]
			}
			if int(dist*100) >= scoreMin {
				ret = true
				break
			}
		}
	}
	return ret
}

func ParseBinFile(vecMap *map[string][DimCnt]float64, binFile string, logger func(string)) int {
	var word string
	var binWordCnt int

	fl, err := os.Open(binFile)
	if err != nil {
		logger(fmt.Sprintf("[ParseBinFile] error=[%v]", err.Error()))
		return 0
	}
	defer fl.Close()

	// 读词总数
	logger("[ParseBinFile] Begin Parsing Bin File...")
	fmt.Fscanf(fl, "%d", &binWordCnt)
	logger(fmt.Sprintf("[ParseBinFile] binWordCnt=[%v]", binWordCnt))

	// 读维度数，为了节约内存，这个参数并没有用到，维度固定为200
	var dim int
	fmt.Fscanf(fl, "%d", &dim)
	if dim != DimCnt {
		logger("[ParseBinFile] Error! dimCnt Do Not Match!")
		return 0
	}

	if *vecMap == nil {
		*vecMap = make(map[string][DimCnt]float64)
	}

	for k, _ := range *vecMap {
		delete(*vecMap, k)
	}

	var vec [DimCnt]float64
	var vecLen float64

	for i := 0; i < binWordCnt; i++ {
		// 读词，这里和C语言版本不一样，空格自动读了
		fmt.Fscanf(fl, "%s", &word)

		// 读向量，二进制转换为float32
		buf := make([]byte, 4*DimCnt)
		f32 := make([]float32, len(buf)/4)
		br := bytes.NewReader(buf)
		fl.Read(buf)
		binary.Read(br, binary.LittleEndian, &f32)

		// 读取换行符，这里和C语言版本不一样，需要显式读1个字节
		c := make([]byte, 1)
		fl.Read(c)

		// 向量计算
		vecLen = 0.0
		for j := 0; j < DimCnt; j++ {
			vec[j] = float64(f32[j])
			vecLen += vec[j] * vec[j]
		}
		vecLen = math.Sqrt(vecLen)
		for j := 0; j < DimCnt; j++ {
			vec[j] /= vecLen
		}

		// 保存到map
		(*vecMap)[word] = vec
	}

	return binWordCnt
}

func ParseClusterFile(clusterMap *map[int][]string, clusterFile string, logger func(string)) int {
	var clusterWordCnt int
	fl, err := os.Open(clusterFile)
	if err != nil {
		logger(fmt.Sprintf("[ParseClusterFile] error=[%v]", err.Error()))
		return 0
	}
	defer fl.Close()
	if *clusterMap == nil {
		*clusterMap = make(map[int][]string)
	}

	for k, _ := range *clusterMap {
		delete(*clusterMap, k)
	}
	logger("[ParseClusterFile] Begin Parsing Cluster File...")
	clusterWordCnt = 0
	for {
		var id string
		var clusterId int
		n, _ := fmt.Fscanf(fl, "%s %d\n", &id, &clusterId)
		if n <= 0 {
			break
		}
		cc, ok := (*clusterMap)[clusterId]
		if ok {
			cc = append(cc, id)
			(*clusterMap)[clusterId] = cc
		} else {
			cc := append(cc, id)
			(*clusterMap)[clusterId] = cc
		}
		clusterWordCnt++
	}
	logger(fmt.Sprintf("[ParseClusterFile] clusterWordCnt=[%v]", clusterWordCnt))
	return clusterWordCnt
}
