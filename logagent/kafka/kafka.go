package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var (
	Client  sarama.SyncProducer
	msgChan chan *sarama.ProducerMessage
)

func Init(address []string, chanSize int64) (err error) {
	// 1.生产者配置
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	// 2.链接kafka
	Client, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		logrus.Error("kafka :producer closed, err:", err)
		return
	}
	msgChan = make(chan *sarama.ProducerMessage, chanSize)
	go sendMsg()
	return err
}
func sendMsg() {
	for {
		select {
		case msg := <-msgChan:
			pid, offset, err := Client.SendMessage(msg)
			if err != nil {
				logrus.Warning("send msgfailed , err", err)
				return
			}
			logrus.Infof("send msg to kafka success, pid :%v, offset:%v", pid, offset)

		}
	}
}

func MsgChan(msg *sarama.ProducerMessage) {
	msgChan <- msg
}
