package taskrunner

import (
	"errors"
	"log"
	"os"
	"streaming-website/conf"
	"sync"

	mydb "streaming-website/db"
)

// task是定制化比较强的部分, runner是可以复用的
// 延迟删除, Dispatcher 从数据库读取要删除的数据,写入 dataChan
// Executor 真正去删除视频文件

// Dispatcher
func VideoClearDispatcher(dc dataChan) error {
	res, err := mydb.ReadVideoDeletionRecord(3)
	if err != nil {
		log.Printf("Video clear dispatcher error: %v", err)
		return err
	}
	if len(res) == 0 {
		return errors.New("all tasks finished")
	}

	// 把 video_del_info 待删除视频id写入 data channel
	for _, id := range res {
		dc <- id
	}

	return nil
}

// Executor
func VideoClearExecutor(dc dataChan) error {
	// 删除视频的goroutine是并行的,需要保存所有的err信息,不知道goroutine的执行快慢
	errMap := &sync.Map{}

FORLOOP:
	for {
		select {
		case vid := <-dc:
			// TODO : NOTICE 并发执行删除视频文件的任务,每个视频文件起1个goroutine
			//  注意 case vid := <-dc: 在循环中 case 的 vid 拿到新的视频id
			//  下面的 go func(){删除视频文件}() 可能没有执行完
			//  Dispatcher 从 video_del_info 表 读取了 1 2 3 写入了 dataChan
			//  Executor 的 goroutine 在后台执行,从 dataChan 也读取到 1 2 3
			//  但是 只删除了 1 此时把 READY_TO_DISPATCH 消息发给了 Dispatcher
			//  导致 Dispatcher 收到 READY_TO_DISPATCH 消息,又去数据库读 但是 2 3 还没删除
			//  Dispatcher 又把 2 3 从 video_del_info 中读出来 ,分发到 dataChan
			//  Executor 有拿到 2 3 导致重复删除数据 但是对业务最终结果没影响 这里就这样写
			go func(video_id interface{}) {
				// 删除视频文件
				err := os.Remove(conf.VideoDir + video_id.(string))
				// 文件不存在的错误 也表示视频被删除了 所以要排除这个错误
				if err != nil && !os.IsNotExist(err) {
					log.Printf("Deleting video error: %v", err)
					// sync.Map 线程安全 多个goroutine并发读写没有问题
					errMap.Store(video_id, err)
					return
				}
				// 删除 video_del_info 表已经删除视频文件id
				if err := mydb.DelVideoDeletionRecord(video_id.(string)); err != nil {
					errMap.Store(video_id, err)
					return
				}
			}(vid)
			// TODO : NOTICE 传参 go func(){}(vid) 这是 goroutine 调用1个闭包
			//  用 goroutine 调用1个闭包函数,只会拿到闭包瞬时的状态,而不会将闭包状态保存
			//  要把参数传给 goroutine
		default:
			break FORLOOP
		}
	}

	var e error
	errMap.Range(func(k, v interface{}) bool {
		e = v.(error)
		// 这里只要有 1个 err存在
		if e != nil {
			return false
		}
		return true
	})

	return e
}
