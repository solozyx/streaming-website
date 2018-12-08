package session

import (
	"sync"
	"db"
	"model"
	"common"
	"time"
)

var(
	// cache 外部Redis mongodb 一个系统增加模块和服务系统复杂度会增加
	// 引入第三方cache带来的好处 小于 开发增量 反而会成为负担
	// 所以本系统没有必要采用Redis session是和用户登录相关的
	// 数据量有限 使用go内置缓存足够
	//
	// sync.Map in_memory缓存 是 go1.9 加入的 在 1.9之前 使用 传统 buildin Map
	// 传统Map非常大的缺点 不支持并发读写 超过2个以上的协程同时并发去读写Map就会panic
	// 使用传统Map要加锁 非常不方便 开发者给社会提了多年意见 增加线程安全的Map
	// 在 1.9加入了 sync.Map 自己实现了线程安全机制 尤其在读上 并发读优化的非常极致
	// 即使 上千并发 上万并发 sync.Map 读取都十分平稳
	// 但是在并发写上 会有问题 并发写 不会报 key conflict 冲突 每次写都加 全局锁 耗时
	// 并发读 并发写 都好的东西 也不容易做到 综合起来 sync.Map 在并发读写表现可以接受
	// 如果用单协程访问 sync.Map 效率的提升并不明显
	sessionMap *sync.Map
)

func init(){
	sessionMap = &sync.Map{}
}

/*
服务启动初始化/重启 从db获取session 到api节点的cache
*/
func LoadSessionsFromDB(){
	var(
		sessionList *sync.Map
		err error
		session *model.SimpleSession
	)
	if sessionList,err = db.RetrieveAllSessions(); err != nil {
		return
	}
	sessionList.Range(func(key, value interface{}) bool {
		// 类型断言
		session = value.(*model.SimpleSession)
		// db session store -> api node cache
		sessionMap.Store(key,session)
		return true
	})
}

/*
用户注册/未登录老用户登录 需要新分配 session_id
*/
func GenerateNewSessionId(login_name string) (session_id string){
	var(
		sessionCreateTime int64
		ttlint64 int64
		session *model.SimpleSession
		err error
	)
	session_id,_ = common.NewUUID()
	// 纳秒 / (1000*1000) 精度是 毫秒
	sessionCreateTime = time.Now().UnixNano() / (1000 * 1000)
	// 精确到毫秒ttl = 当前时间戳 + 服务端session过期时间段 指定30分钟
	ttlint64 = sessionCreateTime + 30 * 60 * 1000
	session = &model.SimpleSession{Username:login_name,TTL:ttlint64}
	if err = db.InsertSession(session_id,ttlint64,login_name); err != nil {
		// session插入数据库成功 才更新 api node session cache
		sessionMap.Store(session_id,session)
	}
	return
}

/*
session过期返回过期状态 未过期返回未过期状态
*/
func IsSessionExpired(session_id string) (login_name string,isSessionExpired bool) {
	var(
		value interface{}
		ok bool
		sessionCreateTime int64
	)
	if value, ok = sessionMap.Load(session_id); ok {
		sessionCreateTime = time.Now().UnixNano() / (1000 * 1000)
		if value.(*model.SimpleSession).TTL < sessionCreateTime {
			deleteExpiredSession(session_id)
			return "", true
		}
		return value.(*model.SimpleSession).Username, false
	}
	// Load失败 认为是session过期
	return "", true
}

func deleteExpiredSession(session_id string){
	sessionMap.Delete(session_id)
	db.DeleteSession(session_id)
}