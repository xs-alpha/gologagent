package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v", err)
		return
	}
	defer cli.Close()

	// watch
	watchCh := cli.Watch(context.Background(), "s4")
	for rsp := range watchCh {
		fmt.Println(rsp)
		for _, evt := range rsp.Events {
			fmt.Printf("evtType:%s, key:%s, value：%s", evt.Type, evt.Kv.Key, evt.Kv.Value)
			//fmt.Println(evt.PrevKv.Value)
		}
	}
}
