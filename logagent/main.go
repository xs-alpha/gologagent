package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"logagent/etcd"
	"logagent/kafka"
	"logagent/tailfile"
)

type Config struct {
	KafkaConfig   `ini:"kafka"`
	CollectConfig `ini:"collect"`
	EtcdConfig    `ini:"etcd"`
}

type EtcdConfig struct {
	Address    string `ini:"address"`
	CollectKey string `ini:"collect_key"`
}

type KafkaConfig struct {
	Address  string `ini:"address"`
	Topic    string `ini:"topic"`
	ChanSize int64  `ini:"chan_size"`
}
type CollectConfig struct {
	LogFilePath string `ini:"logfile_path"`
}

func run() {
	select {}
}

func main() {
	var configObj = new(Config)
	//cfg, err := ini.Load("./conf/config.ini")
	//if err!=nil{
	//	logrus.Error("load config failed , err:%v", err)
	//	return
	//}
	//kafkaAddr := cfg.Section("kafka").Key("address").String()
	//fmt.Println(kafkaAddr)
	err := ini.MapTo(configObj, "./conf/config.ini")
	if err != nil {
		logrus.Errorf("load config failed, err:%v", err)
		return
	}
	fmt.Printf("%#v\n", configObj)
	// 1.初始化链接kafka, 做好准备工作
	err = kafka.Init([]string{configObj.KafkaConfig.Address}, configObj.ChanSize)
	if err != nil {
		logrus.Errorf("init kafka failed, err:%v", err)
		return
	}
	logrus.Info("kafka init success!")
	// 初始化etcd链接
	// 根据配置中的日志路径初始化tail
	err = etcd.Init([]string{configObj.EtcdConfig.Address})
	if err != nil {
		logrus.Errorf("init etcd failed , err:%v", err)
		return
	}
	allConf, err := etcd.GetConf(configObj.EtcdConfig.CollectKey)
	if err != nil {
		logrus.Errorf("etcd get conf failed , err:%v", err)
		return
	}
	fmt.Println("conf:", allConf)
	// 派一个小弟去监控etcd中 configObj.EtcdConfig.CollectKey 对应值的变化
	go etcd.WatchConf(configObj.EtcdConfig.CollectKey)
	// 2.初始化tail
	err = tailfile.Init(allConf)
	if err != nil {
		logrus.Errorf("init tailfile failed , err:%v", err)
		return
	}
	logrus.Info("init tailfile success!")
	run()
	// 3.把日志通过sarama发往kafka
}

// etcdctl put greeting "Hello, etcd"
//bin\windows\kafka-console-consumer.bat --bootstrap-server PLAINTEXT:127.0.0.1:9092 --topic web_log --from-beginning
// bin\windows\zookeeper-server-start.bat config\zookeeper.properties
// bin\windows\kafka-server-start.bat config\server.properties
