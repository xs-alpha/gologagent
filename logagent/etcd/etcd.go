package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"logagent/common"
	"logagent/tailfile"
	"time"
)

var (
	client *clientv3.Client
)

func Init(address []string) (err error) {
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   address,
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		logrus.Error("connect to etcd err:%v", err)
		return err
	}
	return nil
}

func GetConf(key string) (collectEntryList []common.CollectEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	resp, err := client.Get(ctx, key)
	if err != nil {
		logrus.Errorf("get conf from etcd by key: %s failed , err:%v", key, err)
		return
	}
	if len(resp.Kvs) == 0 {
		logrus.Warningf("get len :0 conf from etcd by key:%s,", key)
		return
	}
	ret := resp.Kvs[0]
	fmt.Println("etcd ret:", ret.Value)
	err = json.Unmarshal(ret.Value, &collectEntryList)
	if err != nil {
		logrus.Errorf("jsonUnmarshal failed , err : %v", err)
		return
	}
	return
}

// WatchConf 监控etcd中日志收集项目配置变化的函数
func WatchConf(key string) {
	for {
		watchCh := client.Watch(context.Background(), key)
		// 是覆盖，不用添加
		var newConf []common.CollectEntry
		for wresp := range watchCh {
			for _, evt := range wresp.Events {
				logrus.Info("etcd: get new conf from etcd ", newConf)
				fmt.Printf("etcd: type :%s key:%s value: %s\n", evt.Type, evt.Kv.Key, evt.Kv.Value)
				err := json.Unmarshal(evt.Kv.Value, &newConf)
				if err != nil {
					logrus.Errorf("etcd: json.unmarshal new conf failed, err:%v", err)
					continue
				}
			}
			// 告诉tailfile这个模块应该启用新的配置了
			tailfile.SendNewConf(newConf)
		}
	}
}
