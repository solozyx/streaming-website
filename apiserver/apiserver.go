package main

import (
	"runtime"
	"fmt"
	"flag"
	"time"
	"config"
	"api"
)

var (
	confFileParam string
)

func main() {
	var (
		err error
	)
	initArgs()
	initEnv()
	if err = config.InitConfig(confFileParam); err != nil {
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

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func initArgs(){
	flag.StringVar(&confFileParam,"name","./config.json","传入配置文件路径")
	flag.Parse()
}