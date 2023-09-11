package main

import (
	"flag"
	"fmt"
	"sberTest/pkg/testPack"
	"sync"
)

var numThreads int

func init() {
	flag.IntVar(&numThreads, "threads", 3, "number of threads")
}
func main() {
	flag.Parse()

	var wg sync.WaitGroup

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			testPack.RunTest()
		}()
	}

	wg.Wait()

	fmt.Println("all test are done!")
}
