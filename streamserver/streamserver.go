package main

import (
	"runtime"
	"fmt"
	"time"
	"streamserver"
)

func main() {
	var (
		err error
	)
	initEnv()
	if err = streamserver.InitStreamServer(); err != nil {
		goto ERR
	}
	for{
		time.Sleep(1*time.Second)
	}
	return
ERR:
	fmt.Println(err)
}

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}