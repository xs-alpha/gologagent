package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	d := time.Now().Add(50 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	// 尽管context会过期，但在任何情况下调用他的 cancle函数都是很好的实践
	// 如果不这样， 可能会使其上下文及其父类存活的时间超过必要的时间
	defer cancel()

	select {
	case <-time.After(1 * time.Second):
		fmt.Println("overslept")
	case <-ctx.Done():
		fmt.Println("err:", ctx.Err())
	}
}
