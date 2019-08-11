package main

import (
	"fmt"
	"runtime"
	"time"

	"streaming-website/stream"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)
	if err = stream.InitStreamServer(); err != nil {
		goto ERR
	}
	for{
		time.Sleep(1*time.Second)
	}
	return
ERR:
	fmt.Println(err)
}