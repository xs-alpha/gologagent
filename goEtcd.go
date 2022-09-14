package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		fmt.Print("connect to etcd failed , err:%v, ", err)
		return
	}
	defer cli.Close()

	// put
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//str := `[{"path":"d:/itcast/logs/log.log","topic":"web_log"},{"path":"d:/itcast/001/log/log.log","topic":"web_logs"}]`
	str := `[{"path":"d:/itcast/logs/log.log","topic":"web_log"}]`
	//str := `[{"path":"d:/itcast/logs/log.log","topic":"web_log"},{"path":"d:/itcast/001/log/log.log","topic":"web_logs"},{"path":"d:/itcast/001/log/log2.log","topic":"web_logs2"}]`
	_, err = cli.Put(ctx, "etest", str)
	if err != nil {
		fmt.Printf("put to etcd failed, err%v,", err)
		return
	}
	cancel()

	// get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	get, err := cli.Get(ctx, "etest")
	if err != nil {
		fmt.Printf("get to etcd failed, err%v,", err)
		return
	}
	for _, ev := range get.Kvs {
		fmt.Printf("key :%s, value: %s\n", ev.Key, ev.Value)
	}
	cancel()
	s := []string{"111", "222"}
	for i := range s {
		fmt.Println(i)
	}
}
