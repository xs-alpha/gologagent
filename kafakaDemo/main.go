package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"sync"
)

// kafka消费者
func main() {
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n\n", err)
		return
	}
	partitionList, err := consumer.Partitions("web_log")
	if err != nil {
		fmt.Printf("fail to get list of partition: err%v\n\n", err)
		return
	}
	fmt.Println(partitionList)
	var wg sync.WaitGroup
	for partition := range partitionList {
		pc, err := consumer.ConsumePartition("web_log", int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d, err: %v\n", err)
			return
		}
		defer pc.AsyncClose()
		go func(partitionConsumer sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Println("00000")
				fmt.Printf("Partition:%d offset:%d key:%s Value:%s", msg.Partition, msg.Offset, msg.Key, msg.Value)
			}
		}(pc)
	}
	wg.Wait()
}

// bin\windows\kafka-console-consumer.bat --bootstrap-server localhost:9092 --topic shopping --from-beginning
// bin\windows\zookeeper-server-start.bat config\zookeeper.properties
// bin\windows\kafka-server-start.bat config\server.properties
