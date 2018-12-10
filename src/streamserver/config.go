package streamserver

const(
	// streamserver port
	STREAM_SERFER_PORT = 9000
	// streamserver ReadTimeout
	STREAM_SERFER_READTIMEOUT = 5000
	// streamserver WriteTimeout
	STREAM_SERFER_WRITEIMEOUT = 5000

	// static video file path
	VIDEO_DIR = "./videos/"
	// streamserver concurrent connection limit
	VIDEO_LIMIT = 2
	// static video default type
	VIDEO_DEFAULT_TYPE = "video/mp4"

	// max upload size = 50MB
	MAX_UPLOAD_SIZE = 1024 * 1024 * 50

	// upload file testpage
	PAGE_TEST_UPLOAD_FILE = "upload.html"

	// streamserver internal error
	ERR_INTERNAL = "streamserver internal error"
	// streamserver too many requests
	ERR_TOO_MANY_REQUESTS = "streamserver too many requests"
	// upload video too big
	ERR_UPLOAD_FILR_TOO_BIG = "upload file is too big"

	LOG_SERVER_OPEN_FILE_ERR  = "streamserver server fail to open video file"
	LOG_SERVER_READ_FILE_ERR  = "streamserver fail to read client upload video file"
	LOG_SERVER_WRITE_FILE_ERR = "streamserver server fail to write video file"
	LOG_SERVER_UPLOAD_FILE_SUCCESS = "streamserver upload file success"
)