package main

import (
	"context"
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

func worker(ctx context.Context) {
	defer wg.Done()
Label:
	for {
		fmt.Println("worker...")
		time.Sleep(time.Second)
		select {
		case <-ctx.Done():
			break Label
			// return
		default:

		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go worker(ctx)
	time.Sleep(time.Second * 3)
	cancel()
	wg.Wait()
	fmt.Println("over")
}
