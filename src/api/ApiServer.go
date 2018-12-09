package api

import (
	"net/http"
	"net"
	"time"
	"strconv"
	"config"
	"session"
	"io/ioutil"
	"model"
	"encoding/json"
	"db"
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
		inputBytes []byte
		uc *model.UserCredential
		session_id string
		signUpResp *model.SignedUp
		outputBytes []byte
	)
	// req.Body
	inputBytes,_ = ioutil.ReadAll(req.Body)
	uc = &model.UserCredential{}
	if err = json.Unmarshal(inputBytes,uc); err != nil{
		sendErrorResponse(resp,model.ErrorRequestBodyParseFailed)
		return
	}
	// regiter user to db
	// TODO pwd 加密处理
	if err = db.AddUserCredential(uc.Username,uc.Pwd); err != nil {
		sendErrorResponse(resp,model.ErrorDBError)
		return
	}
	// create session_id to user
	session_id = session.GenerateNewSessionId(uc.Username)
	signUpResp = &model.SignedUp{Success:true,SessionId:session_id}
	if outputBytes,err = json.Marshal(signUpResp); err != nil{
		sendErrorResponse(resp,model.ErrorInternalFaults)
		return
	} else {
		// 201 表示用户注册成功
		sendNormalResponse(resp,outputBytes,201)
	}
}