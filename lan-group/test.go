package main

import (
	"fmt"
	"strconv"
	"sync"
)

var res chan string

var swg sync.WaitGroup

func main1() {
	res = make(chan string, 10)
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go fun(i)
	}
	wg.Wait()
	fmt.Println("close")
	close(res)
	for r := range res {
		fmt.Println(r)
	}
	//for {
	//	if r, ok := <-res; ok {
	//		fmt.Println(r)
	//	} else {
	//		break
	//	}
	//}
}

func fun(a int) {
	defer wg.Done()
	for i := 0; i < 3; i++ {
		res <- strconv.Itoa(i) + " " + strconv.Itoa(a)
	}
}
