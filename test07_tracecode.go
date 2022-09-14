package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type TraceCode string
type UserId string

var wg sync.WaitGroup

var UseridKey = UserId("USERID_KEY")
var TraceCodeKey = TraceCode("TRACE_CODE")

func worker(ctx context.Context) {
	// trace_code, ok := ctx.Value(key).(TraceCode)
	// ctx.Value这是取值
	trace_code, ok := ctx.Value(TraceCodeKey).(string)
	// trace_code, ok := TraceCode("TRACE_CODE").Value(key).(string)
	if !ok {
		fmt.Println("invalid trace  code")
	}
	userid, ok := ctx.Value(UseridKey).(int64)
	if !ok {
		fmt.Println("invalid userid")
	}
	log.Printf("%s worker func ...", trace_code)
	log.Printf("userid : %d worker func ...", userid)
LOOP:
	for {
		fmt.Printf("worker, trace code :%s\n", trace_code)
		time.Sleep(time.Millisecond * 10)
		select {
		case <-ctx.Done():
			break LOOP
		default:
		}
	}
	fmt.Println("worker done！")
	wg.Done()
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	ctx = context.WithValue(ctx, TraceCodeKey, "123456789")
	ctx = context.WithValue(ctx, UseridKey, int64(212121213456))
	log.Printf("%s main 函数", "123456789")
	wg.Add(1)
	go worker(ctx)
	time.Sleep(time.Second * 5)
	cancel()
	wg.Wait()
	fmt.Println("over")
}
