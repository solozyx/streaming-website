package streamserver

import (
	"net/http"
	"net"
	"time"
	"strconv"
	"os"
	"mime/multipart"
	"io/ioutil"
	"log"
	"html/template"
)

var (
	G_streamServer *StreamServer
)

type middleWareHandler struct {
	connLimiter *ConnLimiter
}

func NewMiddleWareHandler(cc int) *middleWareHandler{
	m := &middleWareHandler{
		connLimiter:NewConnLimiter(cc),
	}
	return m
}

type StreamServer struct{
	httpServer *http.Server
	middleWareHandler *middleWareHandler
}

func InitStreamServer() (err error){
	var (
		mux *http.ServeMux
		listener net.Listener
		httpSever *http.Server
		m *middleWareHandler
	)
	mux = http.NewServeMux()
	mux.HandleFunc("/videos",handleStreaming)
	mux.HandleFunc("/upload",handleUpload)
	mux.HandleFunc("/testpage",handleTestPage)
	// bucket token flow limit
	m = NewMiddleWareHandler(VIDEO_LIMIT)
	// listen and serve
	if listener,err = net.Listen("tcp",":" + strconv.Itoa(STREAM_SERFER_PORT)); err != nil {
		return
	}
	httpSever = &http.Server{
		ReadTimeout:time.Duration(STREAM_SERFER_READTIMEOUT) * time.Millisecond,
		WriteTimeout:time.Duration(STREAM_SERFER_WRITEIMEOUT) * time.Millisecond,
		Handler:mux,
	}
	G_streamServer = &StreamServer{
		httpServer:httpSever,
		middleWareHandler:m,
	}
	go httpSever.Serve(listener)
	return
}

/*
多种方式 static video file --> streaming
1.video file -> bit stream 二进制直接发给client 数据量速度 带宽 可控 但是实现复杂
2.web server最简单通用方式 http.ServeContent() 简单网站的视频点播 规模不大都是该方式
*/
func handleStreaming(resp http.ResponseWriter, req *http.Request){
	var(
		err error
		vid string
		videoLink string
		video *os.File
	)

	// bucket token flow limit
	if !G_streamServer.middleWareHandler.connLimiter.GetConn(){
		sendErrorResponse(resp,http.StatusTooManyRequests,ERR_TOO_MANY_REQUESTS)
		return
	}
	defer G_streamServer.middleWareHandler.connLimiter.ReleaseConn()

	// 解析GET表单参数 获取GET参数
	if err = req.ParseForm(); err != nil {
		return
	}
	vid = req.Form.Get("vid")
	videoLink = VIDEO_DIR + vid

	// open static video file
	// TODO 是否需要 buffer io 优化
	if video,err = os.Open(videoLink); err != nil{
		log.Printf(LOG_SERVER_OPEN_FILE_ERR)
		sendErrorResponse(resp,http.StatusInternalServerError,ERR_INTERNAL)
		return
	}
	// 设置response header
	// videoLink如果没有扩展名 真正的video二进制码 是视频mp4格式
	// 这里把文件内容强制设置为mp4 浏览器就能自动以mp4格式解析video 解析完自动组装视频播放
	resp.Header().Set("Content-Type",VIDEO_DEFAULT_TYPE)
	// 播放 ServeContent() 把video的二进制流 传给浏览器 ,浏览器拿到该二进制流
	// 按照 Content-Type 解析 自动播放
	// ServeContent(w ResponseWriter,req *Request,name string,modtime time.Time,
	// content io.ReadSeeker 保证执行ServeContent()时 视频播放流畅 )
	http.ServeContent(resp,req,"",time.Now(),video)
	// 文件指针关闭
	defer video.Close()
}

/*
client static video file --> streaming --> server
*/
func handleUpload(resp http.ResponseWriter, req *http.Request){
	var(
		err error
		file multipart.File
		data []byte
		fileName string
		path string
	)

	// bucket token flow limit
	if !G_streamServer.middleWareHandler.connLimiter.GetConn(){
		sendErrorResponse(resp,http.StatusTooManyRequests,ERR_TOO_MANY_REQUESTS)
		return
	}
	defer G_streamServer.middleWareHandler.connLimiter.ReleaseConn()

	// http.MaxBytesReader 限制 io.Reader 最大读取字节 byte 不是 bit
	req.Body = http.MaxBytesReader(resp,req.Body,MAX_UPLOAD_SIZE)
	// 解析表单
	if err = req.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		sendErrorResponse(resp,http.StatusBadRequest,ERR_UPLOAD_FILR_TOO_BIG)
		return
	}
	// <form name="FormFile key 设置为 file 前端页面需要写好该key ">
	// TODO *multipart.FileHeader 可以验证文件类型 可以在前端页面做验证
	// accept="video/*"
	if file,_,err = req.FormFile("file"); err != nil {
		sendErrorResponse(resp,http.StatusInternalServerError,ERR_INTERNAL)
		return
	}
	if data ,err = ioutil.ReadAll(file); err != nil {
		log.Printf(LOG_SERVER_READ_FILE_ERR)
		sendErrorResponse(resp,http.StatusInternalServerError,ERR_INTERNAL)
		return
	}
	fileName = req.Form.Get("vid")
	path = VIDEO_DIR + fileName
	// 上传文件写入server端磁盘 尽量不用 0777 权限过大 可执行文件就可怕了
	if err = ioutil.WriteFile(path,data,0666); err != nil {
		log.Printf(LOG_SERVER_WRITE_FILE_ERR + " " + path)
		sendErrorResponse(resp,http.StatusInternalServerError,ERR_INTERNAL)
		return
	}
	// send correct response
	log.Printf(LOG_SERVER_UPLOAD_FILE_SUCCESS)
	resp.WriteHeader(http.StatusCreated) // 201
	resp.Write([]byte(LOG_SERVER_UPLOAD_FILE_SUCCESS))
}

func handleTestPage(resp http.ResponseWriter, req *http.Request){
	t,_ := template.ParseFiles(PAGE_TEST_UPLOAD_FILE)
	t.Execute(resp,nil)
}

/*
error response
*/
func sendErrorResponse(resp http.ResponseWriter,statusCode int,errMsg string){
	resp.WriteHeader(statusCode)
	resp.Write([]byte(errMsg))
}