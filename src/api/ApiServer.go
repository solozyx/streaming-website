package api

import (
	"net/http"
	"net"
	"time"
	"strconv"
	"config"
	"common"
)

var (
	G_apiServer *ApiServer
)

type ApiServer struct{
	httpServer *http.Server
}

func InitApiServer() (err error){
	var (
		mux *http.ServeMux
		listener net.Listener
		httpSever *http.Server
	)
	mux = http.NewServeMux()
	mux.HandleFunc(config.G_config.ApiUserRegister,handleUserRegister)
	mux.HandleFunc(config.G_config.ApiUserLogin,handleUserLogin)
	if listener,err = net.Listen("tcp",":" + strconv.Itoa(config.G_config.ApiPort)); err != nil {
		return
	}
	httpSever = &http.Server{
		ReadTimeout:time.Duration(config.G_config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout:time.Duration(config.G_config.ApiWriteTimeout) * time.Millisecond,
		Handler:mux,
	}
	G_apiServer = &ApiServer{httpServer:httpSever}
	go httpSever.Serve(listener)
	return
}

func handleUserRegister(resp http.ResponseWriter, req *http.Request){
	var(
		err error
		bytes []byte
	)
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	if bytes,err = common.BuildResponse(0,"success","success"); err == nil{
		resp.Write(bytes)
	}
	return
ERR:
	if bytes,err = common.BuildResponse(-1,err.Error(),nil); err == nil{
		resp.Write(bytes)
	}
}

func handleUserLogin(resp http.ResponseWriter, req *http.Request){
	var(
		err error
		bytes []byte
		user_name string
	)
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	user_name = req.PostForm.Get("user_name")
	if bytes,err = common.BuildResponse(0,"success",user_name); err == nil{
		resp.Write(bytes)
	}
	return
ERR:
	if bytes,err = common.BuildResponse(-1,err.Error(),nil); err == nil{
		resp.Write(bytes)
	}
}