package taskrunner

/*
Producer / Dispatcher
		|
		|
	controlChan
		|
		|
Consumer / Executor
*/

const(
	// 3种类型dispatch消息
	// d 被调度器收到 下发数据 给 dataChan
	READY_TO_DISPATCH = "d"
	// e 调度器下发数据完成 把e消息发给执行器 执行器从dataChan读取数据 执行
	READY_TO_EXECUTE = "e"
	// 调度器 执行器 任何1个出了问题 或者 没有任务可执行 发送 c消息 把taskrunner常驻任务杀死 结束runner
	CLOSE = "c"
)

type controlChan chan string

// 要下发数据
// interface{} 是go实现泛型的一种机制 可以表示任何类型
type dataChan chan interface{}

// fn 表示1个函数
type fn func(dc dataChan) error