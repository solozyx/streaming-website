package taskrunner

// 核心模块
// runner 会跑一个常驻任务 startDispatcher 等待 runner 的 channel
// control channel 调度器(生产者) 执行器(消费者) 传递数据
// data channel 真正数据的通道

type Runner struct{
	// chan string
	// READY_TO_DISPATCH READY_TO_EXECUTE
	Controller controlChan
	// chan string
	// CLOSE 一般把控制消息 和 错误消息 分开 使用2个 channel
	Error controlChan
	// chan interface{}
	Data dataChan

	dataSize int
	// Runner是否长期存活 true不回收3个channel , false回收3个 channel
	longLived bool

	// 调度器
	// func(dc  chan interface{}) error
	// 是func类型 闭包 go很好实现了闭包功能
	Dispatcher fn
	// 执行器
	// func(dc  chan interface{}) error
	Executor fn
}

// Runner构造方法
func NewRunner(size int,longlived bool,dispatcher fn,executor fn) *Runner {
	return &Runner{
		// buffer = 1 的chan 因为整个过程是非阻塞的
		// 无 buffer chan 可能导致启动Runner卡住
		// 所以一定要用 非阻塞带 buffer 的 chan
		// 放字符串 "d" 和 "e"
		Controller:make(chan string,1),
		// 放字符串 "c"
		Error:make(chan string,1),
		// 每次下发size个数据 1 代表每次只下发1个数据
		Data:make(chan interface{},size),
		longLived:longlived,
		dataSize:size,
		Dispatcher:dispatcher,
		Executor:executor,
	}
}

// 启动scheduler调度任务
func (runner *Runner)StartAll(){
	// Runner初始化完毕 生产者消费者模型 是空的 需要先预置消息
	runner.Controller <- READY_TO_DISPATCH
	runner.startDispatch()
}

// 常驻任务 for 循环 实现消息的异步处理
// 生产者消费者模型
// 通过Data解耦Dispatcher和Executor
// Dispatcher只需要把数据写入Data
// Executor只需要把要执行的任务从Data取出按照任务内容做任务
func (runner *Runner)startDispatch(){
	// 在 startDispatch() 内部声明1个匿名函数闭包并调用
	defer func() {
		if runner.longLived {
			// 常驻则不回收资源
		} else {
			close(runner.Controller)
			close(runner.Error)
			close(runner.Data)
		}
	}()
	// 常驻任务
	for{
		select {
		// 当runner.Controller 有 READY_TO_DISPATCH 和 READY_TO_EXECUTE
		// 调度消息写入 进入 该case
		case c := <- runner.Controller :
			if c == READY_TO_DISPATCH {
				if err := runner.Dispatcher(runner.Data); err != nil {
					// runner.Dispatcher(runner.Data) 调度器
					// 向dataChan写入数据 发生错误 或者 已经写完 都返回 error CLOSE
					// 写入CLOSE 会进入 case e := <- runner.Error: 分支
					runner.Error <- CLOSE
				} else {
					// 调度器把任务下发成功
					runner.Controller <- READY_TO_EXECUTE
				}
			}
			if c == READY_TO_EXECUTE {
				if err := runner.Executor(runner.Data); err != nil {
					runner.Error <- CLOSE
				} else {
					// runner.Executor(runner.Data) 把上1次 Dispatcher 下发给他的数据 取走
					// 并且执行成功 告诉 Dispatcher 进行下1次调度
					// 进入 case c := <- runner.Controller : 分支 下的 if c == READY_TO_DISPATCH 分支
					runner.Controller <- READY_TO_DISPATCH
				}
			}
		// 当有 CLOSE 消息写入 进入 该case
		case e := <- runner.Error:
			if e == CLOSE {
				// 退出 startDispatch()
				return
			}
		default:
			// 上面2个case都不匹配 进入default 分支
			// Controller Error 没有数据 do nothing
		}
	}
}