package stream

type middleWareHandler struct {
	connLimiter *ConnLimiter
}

func NewMiddleWareHandler(cc int) *middleWareHandler {
	m := &middleWareHandler{
		connLimiter: NewConnLimiter(cc),
	}
	return m
}
