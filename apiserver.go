package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"

	"streaming-website/api"
	"streaming-website/conf"
)

var (
	confFileParam string
)

func init()  {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)
	flag.StringVar(&confFileParam,"name","./apiserver.json","传入配置文件路径")
	flag.Parse()
	if err = conf.InitConfig(confFileParam); err != nil {
		goto ERR
	}
	if err = api.InitApiServer(); err != nil {
		goto ERR
	}
	for{
		time.Sleep(1*time.Second)
	}
	return
ERR:
	fmt.Println(err)
}