package api

import (
	"encoding/json"
	"net/http"

	"streaming-website/model"
)

func sendErrorResponse(resp http.ResponseWriter,errResp model.ErrResponse){
	var(
		respJsonBytes []byte
	)
	resp.WriteHeader(errResp.HttpRespStatusCode)
	respJsonBytes,_ = json.Marshal(errResp.Error)
	resp.Write(respJsonBytes)
	// io.WriteString(resp,string(respJsonBytes))
}

func sendNormalResponse(resp http.ResponseWriter,respBytes []byte, statusCode int){
	resp.WriteHeader(statusCode)
	resp.Write(respBytes)
}