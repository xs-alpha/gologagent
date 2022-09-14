package main

import (
	"fmt"
	"github.com/Shopify/sarama"
)

func main() {

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // ack
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 分区
	config.Producer.Return.Successes = true                   // 确认
	// 2.链接kafka
	client, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		fmt.Println("producer closed:", err)
		return
	}
	defer client.Close()

	// 3.封装消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = "shopping"
	msg.Value = sarama.StringEncoder("test消息11")

	// 4.发送消息
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg failed, err:", err)
		return
	}
	fmt.Printf("pid :%v offset:%v\n", pid, offset)
}

// etcdctl put greeting "Hello, etcd"
//bin\windows\kafka-console-consumer.bat --bootstrap-server 127.0.0.1:9092 --topic web_log --from-beginning
// bin\windows\zookeeper-server-start.bat config\zookeeper.properties
// bin\windows\kafka-server-start.bat config\server.properties
// bin\windows\kafka-topics.bat --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic shopping(主题名)
