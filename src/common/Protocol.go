package common

import (
	"encoding/json"
)

type Response struct {
	Errno int `json:"errno"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

func BuildResponse(error int,msg string,data interface{})(resp []byte,err error){
	var(
		response Response
	)
	response.Errno = error
	response.Msg = msg
	response.Data = data
	resp,err = json.Marshal(&response)
	return
}