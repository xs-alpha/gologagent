package tailfile

import (
	"github.com/sirupsen/logrus"
	"logagent/common"
)

type tailTaskMgr struct {
	tailTaskMap      map[string]*tailTask       // 所有的tailTask任务
	collectEntryList []common.CollectEntry      // 所有配置项
	confChan         chan []common.CollectEntry // 等待新配置的通道

}

var (
	ttMgr *tailTaskMgr
)

// Init 为每一个日志文件造一个单独的tailTask
func Init(allconfig []common.CollectEntry) (err error) {
	ttMgr = &tailTaskMgr{
		tailTaskMap:      make(map[string]*tailTask, 20),
		collectEntryList: allconfig,
		confChan:         make(chan []common.CollectEntry),
	}
	// allconfig里存了若干个日志收集项
	// 针对每一个日志收集项创建一个tailObj
	for _, conf := range allconfig {
		// 创建一个日志收集的任务
		tt := newTailTask(conf.Path, conf.Topic)
		err = tt.Init()
		if err != nil {
			logrus.Errorf("create tailobj for path: %s failed, err:%v", conf.Path, err)
			continue
		}
		logrus.Infof("create a tail task for path :%s success", conf.Path)
		ttMgr.tailTaskMap[tt.path] = tt
		// 去收集日志吧
		go tt.run()
	}
	//confChan = make(chan []common.CollectEntry) // 做一个阻塞的channel
	//newConf := <- confChan
	// 新配置来了之后应该管理一下之前启动的那些tailTask
	//logrus.Infof("get new conf from etcd , conf:%v", newConf)
	go ttMgr.watch()
	return
}

func (t *tailTaskMgr) watch() {
	for {
		// 派一个小弟，等着新配置来，
		newConf := <-t.confChan // 取到值说明新的配置来了
		logrus.Infof("start manager tailtask...,  conf :%v", newConf)
		for _, conf := range newConf {
			// 1.原来已经存在的任务不用动
			if t.isExist(conf) {
				continue
			}
			// 2.原来没有的要新创建一个tailTask任务
			tt := newTailTask(conf.Path, conf.Topic)
			err := tt.Init()
			if err != nil {
				logrus.Errorf("create tailobj for path: %s failed, err:%v", conf.Path, err)
				continue
			}
			logrus.Infof("create a tail task for path :%s success", conf.Path)
			// 登记创建的这个tailTask任务，方便后续管理
			t.tailTaskMap[tt.path] = tt
			// 去收集日志吧
			go tt.run()
		}
		// 3.原来有的现在没有的要tailTask停掉
		// TODO:
		for key, task := range t.tailTaskMap {
			var found bool
			for _, conf := range newConf {
				if key == conf.Path {
					found = true
					break
				}
			}
			if !found {
				logrus.Infof("the task collect path:%s need to stop.", task.path)
				// 这个得删了，不然一直有
				delete(t.tailTaskMap, key)
				task.cancel()
			}
		}
	}

}

// 判断tailTask是否存在该收集项
func (t *tailTaskMgr) isExist(conf common.CollectEntry) bool {
	_, ok := t.tailTaskMap[conf.Path]
	return ok
}

func SendNewConf(newConf []common.CollectEntry) {
	ttMgr.confChan <- newConf
}
