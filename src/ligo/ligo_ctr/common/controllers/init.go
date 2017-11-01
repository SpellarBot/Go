package controllers

import (
	"github.com/astaxie/beego"
	"common/utils"
	"common/kafka"
	"ligo/ligo_ctr/common/services"
	"math"
)

var (
	logClient       *utils.LogClient
	robotClient     *utils.DingTalkRobotClient
	kafkaClient     *kafka.KafkaClient
	ctrModel        *services.CTRModel
	status			*string
	modelFile       *string
	password		*string
)

type WebService struct {
	beego.Controller
}

func GlobalInit() error {
	initConf()
	initLogClient()
	initRobotClient()
	initKafkaClient()
	initCTRModel()
	return nil
}




/////////////////////////////////////////////
// 初始化基础服务

func initLogClient() {
	logMaxDays, _ := beego.AppConfig.Int("LogMaxDays")
	logClient = &utils.LogClient{
		LogDir:    beego.AppConfig.String("LogDir"),
		LogPrefix: beego.AppConfig.String("LogPrefix"),
		LogLevel:  beego.AppConfig.String("LogLevel"),
		MaxDays:   logMaxDays,
	}
	logClient.Init()
}



func initRobotClient() {
	dingTalkRobotHost := beego.AppConfig.String("DingTalkRobotHost")
	dingTalkRobotName := beego.AppConfig.String("DingTalkRobotName")
	robotClient = &utils.DingTalkRobotClient{
		RobotHost: dingTalkRobotHost,
		RobotName: dingTalkRobotName,
		Logger:    logClient.Info,
	}
	robotClient.Init()
}

/////////////////////////////////////////////
// 初始化中间件


func initKafkaClient() {
	kafkaClient = &kafka.KafkaClient{
		Topic:         beego.AppConfig.String("KafkaTopic"),
		ConsumerGroup: beego.AppConfig.String("ConsumerGroup"),
		BrokerList:    beego.AppConfig.String("BrokerList"),
		Zookeeper:     beego.AppConfig.String("Zookeeper"),
		Logger:        logClient.Info,
		Debug:         beego.AppConfig.String("LogLevel") == "DEBUG",
	}
	kafkaClient.Init()
}

/////////////////////
///初始化CTRModel
func initCTRModel() {
	ctrModel = &services.CTRModel{}
	var existtime,step,extendtime int64
	var country string
	existtime,_ = beego.AppConfig.Int64("ExistTime")
	step,_  = beego.AppConfig.Int64("Step")
	extendtime,_ = beego.AppConfig.Int64("ExtendTime")
	country = beego.AppConfig.String("Country")
	ctrModel.Init(existtime, step, extendtime, country)
	logClient.Info("Init Status:" + *status)
	if *status == "Restore"{
		er := ctrModel.Restore(*modelFile)
		if er == nil {
			*status = "Work"
			logClient.Info("Restore Success")
		}else {
			*status = "Work"
			logClient.Error("Restore Fail")
		}
	}
}

///初始化开始状态
func initConf(){
	var s1,s2,s3 string
	status = &s1
	modelFile = &s2
	password = &s3
	*status = beego.AppConfig.String("Status")
	*modelFile = beego.AppConfig.String("ModelFile")
	*password  = beego.AppConfig.String("CtrPassword")
}

func CalCtr(show int, play int)(float64){
	var ctr float64
	if show == 0{
		ctr = 0.0
	}else{
		ctr = float64(play)/float64(show)
	}
	ctr = math.Trunc(ctr*1e6)/1e6
	return ctr
}

