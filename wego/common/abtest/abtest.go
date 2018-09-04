package abtest

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/spaolacci/murmur3"
)

const (
	MIN_BUCKETNUM = 100
	DEFAULT_TAG   = "default"
)

type Configuration struct {
	XMLName   xml.Name   `xml:"configuration"`
	Bucketnum int        `xml:"bucketnum"` // 用户总共划分为多少个桶
	Mutex     MutexData  `xml:"mutex"`     // 独占实验
	Layers    LayersData `xml:"layers"`    // 分层实验
}

// 独占实验数据
type MutexData struct {
	XMLName xml.Name    `xml:"mutex"`
	Range   float64     `xml:"range"`   // 独占实验占用实验比例
	Domains DomainsData `xml:"domains"` // 独占实验内部的实验划分
}

// 分层实验数据
type LayersData struct {
	XMLName xml.Name    `xml:"layers"`
	Layer   []LayerData `xml:"layer"`
}

type DomainsData struct {
	XMLName xml.Name     `xml:"domains"`
	Domain  []DomainData `xml:"domain"`
}

type DomainData struct {
	XMLName xml.Name `xml:"domain"`
	Name    string   `xml:"name"`
	Value   float64  `xml:"range"`
}

type LayerData struct {
	XMLName xml.Name    `xml:"layer"`
	Id      string      `xml:"id"`    // 字符串Id
	ExpId   int         `xml:"expid"` // 数值Id
	Domains DomainsData `xml:"domains"`
}

type DomainValue struct {
	Name  string
	Value float64
}

// 保持分流顺序
type ABtestConfigData struct {
	BucketNum    int
	MutexRange   float64
	Mutexs       []DomainValue
	Layers       map[string][]DomainValue
	LayerExpId   map[string]int
	sortedLayers []string
}

func (this *ABtestConfigData) clear() {
	this.BucketNum = 0
	this.MutexRange = 0
	this.Mutexs = this.Mutexs[:0]
	for k, _ := range this.Layers {
		delete(this.Layers, k)
		delete(this.LayerExpId, k)
	}
}

type ABtestConfig struct {
	Logger func(string)

	globalIndex int
	config1     ABtestConfigData
	config2     ABtestConfigData
}

// 初始化
func (this *ABtestConfig) Init(filePath string, period int) error {
	if this.Logger == nil {
		this.Logger = func(msg string) {
			fmt.Println(msg)
		}
	}
	this.globalIndex = 0
	this.GetConfAndSwitch(filePath)
	err := this.checkEqualsOne()
	if err != nil {
		return err
	}
	if period > 0 {
		go func(period int) {
			ticker := time.NewTicker(time.Duration(period) * time.Minute)
			for _ = range ticker.C {
				this.GetConfAndSwitch(filePath)
			}
		}(period)
	}
	return nil
}

func (this *ABtestConfig) Print() {
	if this.globalIndex == 1 {
		this.printABTestConfigData(this.config1)
	} else {
		this.printABTestConfigData(this.config2)
	}
}

func (this *ABtestConfig) printABTestConfigData(a ABtestConfigData) {
	this.Logger("=====================ABTestConfig=====================")
	this.Logger(fmt.Sprintf("%15s : %d", "BucketNum", a.BucketNum))
	this.Logger(fmt.Sprintf("%15s : %1.3f", "MutextRange", a.MutexRange))
	this.Logger("---------------------Mutex---------------------")
	for k, Mutex := range a.Mutexs {
		this.Logger(fmt.Sprintf("    %d Mutex : %15s %1.3f", k, Mutex.Name, Mutex.Value))
	}
	this.Logger("---------------------Layer---------------------")
	for LayerName, Layer := range a.Layers {
		this.Logger(fmt.Sprintf("    %d Layer : %s", a.LayerExpId[LayerName], LayerName))
		for _, Value := range Layer {
			this.Logger(fmt.Sprintf("        %15s %1.3f", Value.Name, Value.Value))
		}
	}
	this.Logger("=====================ABTestConfig=====================")
}

// 本地XML文件加载配置文件
func (this *ABtestConfig) GetConfAndSwitch(filePath string) {
	this.Logger("ABtestConfig Load Config File")
	reader, err := os.Open(filePath)
	if err != nil {
		this.Logger(fmt.Sprintf("ABtestConfig Conf Switch err: %v", err))
		return
	}
	defer reader.Close()
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		this.Logger(fmt.Sprintf("ABtestConfig Conf Switch err: %v", err))
		return
	}

	conf, ok := this.parseConf(data)
	if !ok {
		return
	}

	abtestConfigData := this.getAnotherABtestConf()
	abtestConfigData.clear()
	abtestConfigData.BucketNum = conf.Bucketnum
	abtestConfigData.MutexRange = conf.Mutex.Range
	for _, domain := range conf.Mutex.Domains.Domain {
		if abtestConfigData.Mutexs == nil {
			abtestConfigData.Mutexs = []DomainValue{}
		}
		abtestConfigData.Mutexs = append(abtestConfigData.Mutexs, DomainValue{domain.Name, domain.Value})
	}
	if abtestConfigData.Layers == nil {
		abtestConfigData.Layers = make(map[string][]DomainValue)
	}
	if abtestConfigData.LayerExpId == nil {
		abtestConfigData.LayerExpId = make(map[string]int)
	}
	for _, layer := range conf.Layers.Layer {
		if _, ok := abtestConfigData.Layers[layer.Id]; !ok {
			abtestConfigData.Layers[layer.Id] = []DomainValue{}
			abtestConfigData.LayerExpId[layer.Id] = layer.ExpId
		}
		for _, domain := range layer.Domains.Domain {
			abtestConfigData.Layers[layer.Id] = append(abtestConfigData.Layers[layer.Id], DomainValue{domain.Name, domain.Value})
		}
	}
	if this.globalIndex == 0 || this.globalIndex == 2 {
		this.globalIndex = 1
	} else {
		this.globalIndex = 2
	}
	abtestConfigData.setSortedLayers()
}

// 以[]byte内容转入，以支持其它配置方式
func (this *ABtestConfig) parseConf(data []byte) (Configuration, bool) {
	v := Configuration{}
	err := xml.Unmarshal(data, &v)
	if err != nil {
		this.Logger(fmt.Sprintf(" ABtestConfig parseConf err: %v", err))
		return v, false
	}
	return v, true
}

// 获取当前在用的ABTestConfig
func (this *ABtestConfig) getCurrentABtestConf() *ABtestConfigData {
	if this.globalIndex == 1 {
		return &this.config1
	} else {
		return &this.config2
	}
}

// 获取当前不在用的ABTestConfig
func (this *ABtestConfig) getAnotherABtestConf() *ABtestConfigData {
	if this.globalIndex == 1 {
		return &this.config2
	} else {
		return &this.config1
	}
}

// 实验比例之和=1校验
func (this *ABtestConfig) checkEqualsOne() error {
	conf := this.getCurrentABtestConf()
	for layerId, layer := range conf.Layers {
		sum := 0.0
		for _, domain := range layer {
			sum += domain.Value
		}
		if !(math.Dim(sum, 1.0) >= 0 && math.Dim(1.0, sum) >= 0) {
			this.Logger(fmt.Sprintf("Check abtest.xml Not Equals 1.0! layerId=%s", layerId))
			return errors.New("Check abtest.xml Not Equals 1.0")
		}
	}
	if this.getCurrentABtestConf().BucketNum <= MIN_BUCKETNUM {
		return errors.New("BucketNum Is Too Small")
	}
	return nil
}

// hash函数：一致性hash算法
func (this *ABtestConfig) hash(data []byte) uint64 {
	h64Byte := murmur3.New64()
	h64Byte.Write(data)
	hash := h64Byte.Sum64()
	return hash
}

// 根据userId获取用户落在Mutex还是Layers，并返回相应的桶
func (this *ABtestConfig) GetTag(userid string) (string, string, bool) {
	conf := this.getCurrentABtestConf()
	bucketNum := conf.BucketNum
	mutexRange := conf.MutexRange
	userHash := this.hash([]byte(userid))
	userBucket := userHash % uint64(bucketNum)
	if userBucket < uint64(float64(bucketNum)*mutexRange) {
		currentRange := 0.0
		userHash = this.hash([]byte("mutex" + strconv.FormatUint(userBucket, 10)))
		userBucket := userHash % uint64(bucketNum)
		for _, domain := range conf.Mutexs {
			if userBucket < uint64(float64(bucketNum)*(currentRange+domain.Value)) {
				return "mutex", domain.Name, true
			} else {
				currentRange += domain.Value
			}
		}
	} else {
		layers := ""
		k := 0

		for _, layerId := range conf.sortedLayers {
			layerName := this.getLayeTag(layerId, userHash)
			k++
			if k == 1 {
				layers = fmt.Sprintf("%s_%s", layerId, layerName)
			} else {
				layers = layers + "," + fmt.Sprintf("%s_%s", layerId, layerName)
			}
		}
		return "layer", layers, true
	}
	return "", "", false
}

// 根据用户Hash值获取用户在Layer某一层的Abtag
func (this *ABtestConfig) getLayeTag(layerId string, userBucket uint64) string {
	conf := this.getCurrentABtestConf()
	Layers := conf.Layers
	bucketNum := conf.BucketNum
	userHash := this.hash([]byte("layer" + layerId + strconv.FormatUint(userBucket, 10)))
	userBucket = userHash % uint64(bucketNum)
	currentRange := 0.0
	layer := Layers[layerId]
	for _, domain := range layer {
		if userBucket < uint64(float64(bucketNum)*(currentRange+domain.Value)) {
			return domain.Name
		} else {
			currentRange += domain.Value
		}
	}
	return DEFAULT_TAG
}

// 获取按实验ID顺序的Layer层
func (this *ABtestConfigData) setSortedLayers() {
	sortedLayers := make([]string, 0)
	expIds := make([]int, 0)
	layersIndex := make(map[int]string)
	for k, v := range this.LayerExpId {
		expIds = append(expIds, v)
		layersIndex[v] = k
	}
	sort.Ints(expIds)
	for _, index := range expIds {
		sortedLayers = append(sortedLayers, layersIndex[index])
	}
	this.sortedLayers = sortedLayers
}
