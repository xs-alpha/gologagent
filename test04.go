package main

import (
	"context"
	"fmt"
)

// gen 返回的是一个只读的channel，channel是int类型的
func gen(ctx context.Context) <-chan int {
	dst := make(chan int)
	n := 1
	// 这个函数一直往通道里面放值
	go func() {
		for {
			select {
			// 从通道里面取值可以不接受
			case <-ctx.Done():
				// return 结束该goroutine 防止泄露
				return
			case dst <- n:
				fmt.Println(n)
				n++
			}
		}
	}()
	return dst
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	// 当取完需要的整数后调用cancel
	defer cancel()

	for n := range gen(ctx) {
		fmt.Println(n)
		if n == 5 {
			break
		}
	}
}
