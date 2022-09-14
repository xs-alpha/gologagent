package main

import (
	"fmt"
	"sync"
	"time"
)

// 补充知识点，
// make, new区别
// 都是用来初始化内存
// new 多用来为基本数据类型（bool, string, int...）初始化内存，返回的是指针
// make 用来初始化slice, map, channel，返回的是对应类型

var wg sync.WaitGroup

func worker(ch <-chan bool) {
	defer wg.Done()
Label:
	for {
		fmt.Println("worker...")
		time.Sleep(time.Second)
		select {
		case <-ch:
			break Label
		default:

		}
	}
}

func main() {
	var exitChan = make(chan bool)
	wg.Add(1)
	go worker(exitChan)
	time.Sleep(time.Second * 3)
	exitChan <- true
	wg.Wait()
	fmt.Println("over")
}
