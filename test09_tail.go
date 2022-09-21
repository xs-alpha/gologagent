package main

import (
	"fmt"
	"github.com/hpcloud/tail"
	"time"
)

func main() {
	fileName := `./xx.log`
	config := tail.Config{
		ReOpen:    true, // 文件超过大小分割以后，自动打开
		Follow:    true, // 自动跟随文件名，
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	tails, err := tail.TailFile(fileName, config)
	if err != nil {
		fmt.Println("test09- tailfile %s failed, err:%v \n", fileName, err)
		return
	}
	var (
		msg *tail.Line
		ok  bool
	)
	for {
		msg, ok = <-tails.Lines // chan tailfile.Line
		if !ok {
			fmt.Printf("tailfile file close reopen, filename:%s \n")
			time.Sleep(time.Second) // 读取出错, 等待一秒
			continue
		}
		fmt.Println("msg:", msg.Text)
	}
}
