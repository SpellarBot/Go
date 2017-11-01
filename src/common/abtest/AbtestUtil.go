package abtest

import (
	"encoding/xml"
	"os"
	"io/ioutil"
	"bytes"
	"github.com/spaolacci/murmur3"
	"strconv"
	"fmt"
	"time"
)

const DEFAULT_TAG  = "default"

type Configuration struct {
	XMLName 	xml.Name `xml:"configuration"`
	Bucketnum	int `xml:"bucketnum"`
	Mutex     MutexData `xml:"mutex"`
	Layers     LayersData `xml:"layers"`
}

type MutexData struct {
	XMLName  xml.Name `xml:mutex`
	Range    float64 `xml:"range"`
	Domains  DomainsData `xml:"domains"`
}

type DomainsData struct {
	XMLName  xml.Name `xml:"domains"`
	Domain   []DomainData `xml:"domain"`
}

type DomainData struct {
	XMLName    xml.Name `xml:"domain"`
	Name  string `xml:"name"`
	Value  float64	`xml:"range"`
}

type LayersData struct {
	XMLName  xml.Name `xml:"layers"`
	Layer  []LayerData `xml:"layer"`
}

type LayerData struct {
	XMLName  xml.Name `xml:"layer"`
	Id 	string `xml:"id"`
	Domains  DomainsData `xml:"domains"`
}

type DomainValue struct {
	Name   string
	Value  float64
}

// 保持分流顺序
type ABtestConfigData struct {
	BucketNum	int
	MutexRange	float64
	Mutexs	[]DomainValue // domain : range
	Layers	map[string][]DomainValue // layerid : domain : range
}

func (this *ABtestConfigData) clear()  {
	this.BucketNum = 0
	this.MutexRange = 0
	this.Mutexs = this.Mutexs[:0]
	for k, _ := range this.Layers {
		delete(this.Layers, k)
	}
}

type ABtestConfig struct {
	Logger    func(string)

	globalIndex int // 0:init, 1:map1, 2:map2
	config1   ABtestConfigData
	config2   ABtestConfigData
}

// 初始化, 使用之前必须先调用此方法
func (this *ABtestConfig) Init(filePath string, period int)  {
	if this.Logger == nil {
		this.Logger = func(msg string) {
			fmt.Println(msg)
		}
	}
	this.globalIndex = 0
	this.GetConfAndSwitch(filePath)
	// 周期加载，以支持动态加载ABtest配置文件，实现动态切换流量
	if period > 0 {
		go func(period int) {
			ticker := time.NewTicker(time.Duration(period) * time.Minute)
			for _ = range ticker.C {
				this.GetConfAndSwitch(filePath)
			}
		}(period)
	}
}

/**
当前只支持从本地文件系统加载配置文件，以后有需求再增加其它方式加载，可以直接调用此方法，实现API动态触发切换
 */
func (this *ABtestConfig) GetConfAndSwitch(filePath string) {
	this.Logger("[ABtestConfig][GetConfAndSwitch] load config file.")
	reader, err := os.Open(filePath)
	if err != nil {
		this.Logger(fmt.Sprintf("[ABtestConfig][GetConfAndSwitch] err=[%v]", err))
		return
	}
	defer reader.Close()
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		this.Logger(fmt.Sprintf("[ABtestConfig][GetConfAndSwitch] err=[%v]", err))
		return
	}

	conf, ok:= this.parseConf(data)
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
	for _, layer := range conf.Layers.Layer {
		if abtestConfigData.Layers == nil {
			abtestConfigData.Layers = make(map[string][]DomainValue)
		}
		for _, domain := range layer.Domains.Domain {
			if _, ok := abtestConfigData.Layers[layer.Id]; !ok {
				domains := []DomainValue{}
				abtestConfigData.Layers[layer.Id] = domains
			}
			abtestConfigData.Layers[layer.Id] = append(abtestConfigData.Layers[layer.Id], DomainValue{domain.Name, domain.Value})
		}
	}

	// 索引切换，更新模型数据
	if this.globalIndex == 0 || this.globalIndex == 2 {
		this.globalIndex = 1
	} else {
		this.globalIndex = 2
	}

}
/**
以[]byte内容转入，以支持其它配置方式
 */
func (this *ABtestConfig) parseConf(data []byte) (Configuration, bool) {
	v := Configuration{}
	err := xml.Unmarshal(data, &v)
	if err != nil {
		this.Logger(fmt.Sprintf("[ABtestConfig][parseConf] err=[%v]", err))
		return v, false
	}
	return v, true
}

func (this *ABtestConfig) getABTagImpl(userId string, layerId string, option string) (string, bool) {
	if userId == "" {
		return "", true
	}
	sb := bytes.Buffer{}
	sb.WriteString(userId)
	if layerId != "" {
		sb.WriteString(layerId)
	}
	if option != "" {
		sb.WriteString(option)
	}
	userHash := this.getHash([]byte(sb.Bytes()))

	return this.getTag(layerId, userHash)
}

func (this *ABtestConfig) getHash(data []byte) uint64 {
	h64Byte := murmur3.New64()
	h64Byte.Write(data)
	hash := h64Byte.Sum64()
	return hash
}

func (this *ABtestConfig) getTag(layerId string, userHash uint64) (string, bool) {
	conf := this.getCurrentABtestConf()
	bucketNum := conf.BucketNum
	mutexRange := conf.MutexRange

	// 如查没有初始完成，统一返回DEFAULT_TAG
	if bucketNum == 0 {
		return "", false
	}

	userBucket := userHash % uint64(bucketNum)// + uint64(bucketNum / 2)

	if userBucket < uint64(float64(bucketNum) * mutexRange) { // 落在独占区
		sb := bytes.Buffer{}
		sb.WriteString("mutex")
		sb.WriteString(strconv.FormatUint(userBucket, 10))
		newHash := this.getHash(sb.Bytes())

		userBucket = newHash % uint64(bucketNum)

		if layerId != "" || conf.Mutexs == nil { // 独占区layerId应该为空
			return "", false
		} else {
			currentRange := 0.0
			for _, domain := range conf.Mutexs {
				if userBucket < uint64(float64(bucketNum) * (currentRange + domain.Value)) {
					return domain.Name, true;
				} else {
					currentRange += domain.Value
				}
			}
		}
	} else { // 在并行区
		if mutexRange > 0 {
			sb := bytes.Buffer{}
			sb.WriteString(strconv.FormatUint(userBucket, 10))
			sb.WriteString("layers")
			newHash := this.getHash(sb.Bytes())

			userBucket = newHash % uint64(bucketNum)
		}
		if conf.Layers == nil {
			this.Logger("[ABtestConfig][getTag] conf.Layers == nil")
			return "", false
		}

		layer, ok := conf.Layers[layerId]
		if layerId == "" || !ok || layer == nil{
			return "", false
		}

		currentRange := 0.0
		for _, domain := range layer {
			if userBucket < uint64(float64(bucketNum) * (currentRange + domain.Value)) {
				return domain.Name, true;
			} else {
				currentRange += domain.Value
			}
		}
	}
	return DEFAULT_TAG, false
}

/**
获取user在并行区相应layer的ABtag，实验中，必须实现'default' 这个ABtag的实验分支的代码,
option 可以为空
return: tag, iswant . iswant=false， 表明这是排除在外的流量
 */
func (this *ABtestConfig) GetABTag(userId string, layerId string) (string, bool)  {
	return this.getABTagImpl(userId, layerId, "")
}
func (this *ABtestConfig) GetABTag2(userId string, layerId string, option string) (string, bool)  {
	return this.getABTagImpl(userId, layerId, option)
}

/**
获取user在独占区相应layer的ABtag，实验中，必须实现'default' 这个ABtag的实验分支的代码
return: tag, iswant . iswant=false， 表明这是排除在外的流量
 */
func (this *ABtestConfig) GetMutexABTag(userId string) (string, bool) {
	return this.getABTagImpl(userId, "", "")
}

func (this *ABtestConfig) getCurrentABtestConf() *ABtestConfigData {
	if this.globalIndex == 1 {
		return &this.config1
	} else {
		return &this.config2
	}
}

func (this *ABtestConfig) getAnotherABtestConf() *ABtestConfigData {
	if this.globalIndex == 1 {
		return &this.config2
	} else {
		return &this.config1
	}
}