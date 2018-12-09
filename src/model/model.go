package model

/*
用户
*/
type UserCredential struct {
	Username string `json:"login_name"`
	Pwd string `json:"pwd"`
}

type SignedUp struct {
	Success bool `json:"success"`
	SessionId string `json:"session_id"`
}

/*
video
*/
type VideoInfo struct{
	// 视频Id
	Id string
	// 作者Id
	AuthorId int
	// 视频名称
	Name string
	// 显示创建时间
	DisplayCtime string
	// 在逻辑层意义不大 只在database操作中有意义 不在model定义
	// CreateTime
}

/*
comment
*/
type Comment struct{
	// 评论Id
	Id string
	// 视频Id
	VideoId string
	// 作者名
	Author string
	// 评论内容
	Content string
}

/*
session 是否过期 存储
1.服务启动初始化/重启 从db获取session 到api节点的cache
2.新用户注册/未登录老用户登录 需要新分配 session_id
3.session过期返回过期状态 未过期返回未过期状态

NOTICE 没有修改session功能 一个session写入了有固定的session_id 和 TTL
修改会对业务复杂度造成冲击 所以session不支持修改
*/
type SimpleSession struct{
	// users.login_name
	Username string
	// 校验session是否过期
	TTL int64
}

type Err struct{
	Error string `json:"error"`
	ErrorCode string `json:"error_code"`
}

type ErrResponse struct{
	HttpRespStatusCode int
	Error Err
}

var(
	ErrorRequestBodyParseFailed = ErrResponse{
		HttpRespStatusCode:400,
		Error:Err{
			Error:"request body is not correct",
			ErrorCode:"001",
		},
	}

	ErrorNotAuthUser = ErrResponse{
		HttpRespStatusCode:401,
		Error:Err{
			Error:"user authentication failed",
			ErrorCode:"002",
		},
	}

	ErrorDBError = ErrResponse{
		HttpRespStatusCode:500, // 数据库错误是 500 内部错误
		Error:Err{
			Error:"DB error",
			ErrorCode:"003",
		},
	}

	ErrorInternalFaults = ErrResponse{
		HttpRespStatusCode:500, // json序列化反序列化错误 500 内部错误
		Error:Err{
			Error:"Internal server error",
			ErrorCode:"004",
		},
	}
)