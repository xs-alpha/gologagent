package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

// 代码连接etcd

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

	// put
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//str := `[{"path":"d:/logs/s4.log","topic":"s4_log"},{"path":"e:/logs/web.log","topic":"web_log"}]`
	//str := `[{"path":"d:/logs/s4.log","topic":"s4_log"},{"path":"e:/logs/web.log","topic":"web_log"},{"path":"c:/logs/nazha.log","topic":"nazha_log"}]`
	str := `[{"path":"d:/logs/s4.log","topic":"s4_log"},{"path":"e:/logs/web.log","topic":"web_log"},{"path":"c:/logs/nazha.log","topic":"nazha_log"},{"path":"c:/logs/nazha3.log","topic":"222_log"}]`
	_, err = cli.Put(ctx, "collect_log_conf", str)
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v", err)
		return
	}
	cancel()

	// get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	gr, err := cli.Get(ctx, "collect_log_conf")
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v", err)
		return
	}
	for _, ev := range gr.Kvs {
		fmt.Printf("key:%s value:%s\n", ev.Key, ev.Value)
	}
	cancel()
}
