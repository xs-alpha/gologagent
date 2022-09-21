package main

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
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
	str := `[{"path":"/home/xiaosheng/go/tmp/tmp1/log1.log","topic":"web_log"},{"path":"/home/xiaosheng/go/tmp/tmp2/log2.log","topic":"web_logs"}]`
	// str := `[{"path":"/home/xiaosheng/go/tmp/tmp1/log1.log","topic":"web_log"}]`
	//str := `[{"path":"/home/xiaosheng/go/tmp/tmp1/log1.log","topic":"web_log"},{"path":"/home/xiaosheng/go/tmp/tmp2/log2.log","topic":"web_logs"},{"path":"/home/xiaosheng/go/tmp/tmp3/log3.log","topic":"web_logs2"}]`
	_, err = cli.Put(ctx, "etest", str)
	if err != nil {
		fmt.Printf("put to etcd failed, err%v,", err)
		return
	}
	cancel()

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

// x:="- /home/xiaosheng/go/tmp/tmp1/log1.log"
// get
