package tailfile

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"logagent/common"
	"logagent/kafka"
	"strings"
	"time"
)

var (
	TailObj  *tail.Tail
	confChan chan []common.CollectEntry
)

type tailTask struct {
	path   string
	topic  string
	tobj   *tail.Tail
	ctx    context.Context
	cancel context.CancelFunc
}

// newTailTask 根据topic 和path造一个tailTask对象
func newTailTask(path, topic string) *tailTask {
	ctx, cancel := context.WithCancel(context.Background())
	tt := &tailTask{
		path:   path,
		topic:  topic,
		ctx:    ctx,
		cancel: cancel,
	}
	return tt
}

// run 真正读日志往kafka里发送的方法
func (t *tailTask) run() {
	logrus.Infof("collect for path:%s is running...", t.path)
	// 读取日志，发往kafaka
	for {
		select {
		case <-t.ctx.Done():
			logrus.Infof("path %s is stoping...", t.path)
			return
		case msg, ok := <-t.tobj.Lines:
			if !ok {
				logrus.Warnf("tail file close reopen, path:%s\n", t.path)
				time.Sleep(time.Second)
				continue
			}
			if len(strings.Trim(msg.Text, "\r")) == 0 {
				logrus.Infof("出现空行，跳过")
				continue
			}
			logrus.Infof("msg:%v", msg.Text)
			msg2 := &sarama.ProducerMessage{}
			msg2.Topic = t.topic
			msg2.Value = sarama.StringEncoder(msg.Text)
			kafka.MsgChan(msg2)
		}
	}
}

// Init 使用tail包打开真正的日志文件准备读
func (t *tailTask) Init() (err error) {
	config := tail.Config{
		ReOpen:    true,
		MustExist: false,
		Poll:      true,
		Follow:    true,
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: 2,
		},
	}
	t.tobj, err = tail.TailFile(t.path, config)
	return
}
