package streamserver

import "log"

/*
流控模块 限制连接数
实现非常小 但是起到非常大的作用
go buffer channel 在并发处理先天优势
*/
type ConnLimiter struct{
	// bucket token 并发容量
	concurrentConn int
	// bucket 桶
	bucket chan int
}

/*
go没有原生的构造函数 每个struct初始化需要自己写构造函数
*/
func NewConnLimiter(cc int) (*ConnLimiter) {
	return &ConnLimiter{
		concurrentConn:cc,
		bucket:make(chan int,cc),
	}
}

/*
request进来时刻调用
*/
func (connLimiter *ConnLimiter)GetConn() bool {
	if len(connLimiter.bucket) >= connLimiter.concurrentConn {
		log.Printf("streamserver 流控模块 reached the rate limitation ")
		return false
	}
	// 随便写入1个int值 一般写 1 表示从令牌桶获取到1个token
	connLimiter.bucket <- 1
	return true
}

/*
request完成对应的response 回收 token
*/
func (connLimiter *ConnLimiter)ReleaseConn(){
	token := <- connLimiter.bucket
	log.Printf("streamserver 流控模块 new conn coming: %d",token)
}