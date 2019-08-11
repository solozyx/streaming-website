package taskrunner

import (
	"errors"
	"log"
	"testing"
	"time"
)

func TestRunner(t *testing.T) {
	dispatcher := func(dc dataChan) error {
		for i := 0; i < 30; i++ {
			dc <- i
			log.Printf("Dispatcher send: %v", i)
		}
		return nil
	}

	executor := func(dc dataChan) error {
	FORLOOP:
		for {
			select {
			case data := <-dc:
				log.Printf("Executor receive: %v", data)
			default:
				break FORLOOP
			}
		}
		// return nil
		// executor执行完毕 就没有任务了
		return errors.New("Executor executed all jobs")
	}

	runner := NewRunner(30, false, dispatcher, executor)
	// runner.StartAll() 不能直接调用 要起1个goroutine执行
	// 因为 func (runner *Runner)startDispatch(){ for{ 这里死循环 } }
	// 如果后台不用 goroutine 则无法执行 下面的 sleep操作
	go runner.StartAll()
	time.Sleep(3 * time.Second)
}
