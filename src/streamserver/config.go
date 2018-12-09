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

	// streamserver internal error
	ERR_INTERNAL = "streamserver internal error"
	// streamserver too many requests
	ERR_TOO_MANY_REQUESTS = "streamserver too many requests"
)