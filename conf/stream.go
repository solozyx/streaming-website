package conf

const(
	StreamServerPort = 9000
	StreamServerReadTimeout = 5000
	StreamServerWriteTimeout = 5000

	VideoDir = "./videos/"

	VideoLimit = 2

	VideoDefaultType = "video/mp4"

	// max upload size = 50MB
	MaxUploadSize = 1024 * 1024 * 50
)