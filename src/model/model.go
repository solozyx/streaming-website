package model

/*
用户
*/
type UserCredential struct {
	Username string `json:"user_name"`
	Pwd string `json:"pwd"`
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