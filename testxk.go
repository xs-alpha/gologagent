package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	lis := []string{"119", "179", "179", "119"}
	t := -1
	r := -1
	for i := 0; i < 500000; i++ {
		rand.Seed(time.Now().UnixNano())
		i := rand.Intn(4)
		fmt.Println("oxoo", i)
		if i == 0 || i == 3 {
			t++
		} else {
			r++
		}
	}
	if t > r {
		fmt.Println(lis[0])
	} else if t < r {
		fmt.Println(lis[1])
	} else {
		fmt.Println("another")
	}
}
