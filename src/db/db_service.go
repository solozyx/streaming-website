package db

import (
	"database/sql"
	"log"
	"model"
	"common"
	"time"
)

// 时间格式化固定字符串 任何一个字符都不能改动
var goParseTimeFormat = "Jan 02 2006, 15:04:05"

/*
添加用户
*/
func AddUserCredential(loginName string,pwd string) (err error) {
	var(
		stmtIns *sql.Stmt
	)
	var insertSql = "INSERT INTO users (login_name,pwd) VALUES (?,?)"
	if stmtIns,err = dbConn.Prepare(insertSql);err != nil{
		return
	}
	if _,err = stmtIns.Exec(loginName,pwd); err != nil {
		// 这里return后 stmtIns.Close() 执行不到 导致对象越开越多 需要使用 defer
		return
	}
	// 提高性能 这里 不写 defer stmtIns.Close()
	defer stmtIns.Close()
	return
}

/*
获取用户信息
*/
func GetUserCredential(loginName string) (pwd string,err error) {
	var(
		stmtOut *sql.Stmt
	)
	var selectSql = "SELECT pwd FROM users WHERE login_name = ?"
	if stmtOut,err = dbConn.Prepare(selectSql);err != nil{
		log.Printf("mysql 查询 user 表 pwd 字段失败 %s",err.Error())
		return "",err
	}
	// query 1 row data
	if err = stmtOut.QueryRow(loginName).Scan(&pwd); err != nil {
		if err != sql.ErrNoRows {
			// sql.ErrNoRows 当 stmtOut.QueryRow 返回的查询结果是 nil空值 0 row
			// Scan 会把这个错误带出来 no rows
			// no rows 不是真正的处理错误 而是数据库本身就没有这条记录
			return "",err
		}
		if err == sql.ErrNoRows {
			return "",nil
		}
	}
	defer stmtOut.Close()
	return pwd,nil
}

/*
删除用户
*/
func DeleteUser(loginName string,pwd string) (err error){
	var(
		stmtDel *sql.Stmt
	)
	var deleteSql = "DELETE FROM users WHERE login_name=? and pwd=?"
	if stmtDel,err = dbConn.Prepare(deleteSql); err != nil {
		log.Printf("mysql 删除 user 表记录失败 %s",err.Error())
		return
	}
	if _,err = stmtDel.Exec(loginName,pwd); err != nil{
		return
	}
	defer stmtDel.Close()
	return
}

/*
添加video
@param  aid author id int
@param  name 视频名称 string
@return videoInfo object
*/
func AddNewVideo(aid int,name string) (videoInfo *model.VideoInfo,err error){
	var(
		vid string
		t time.Time
		ctime string
		stmtIns *sql.Stmt
	)
	// Create uuid
	if vid,err = common.NewUUID(); err != nil{
		return nil,err
	}
	// createtime 是进入到 AddNewVideo 函数的时间
	// video_info.create_time 字段是记录插入数据库的时间
	// createtime 和 video_info.create_time 这2个时间是相似的
	// 写数据库在 1秒之内 或者 几十毫秒之内
	// video_info.create_time 可以按入库时间排序
	// 显示在前端页面的时间是 displayCtime C create
	// TODO NOTICE
	// 这2个时间会否错乱?
	// 进入 AddNewVideo 函数 是在1个goroutine下 整个操作流程是顺序串行的
	// 多个video写库它们的 displayCtime 相对的顺序和对应的video_info.create_time
	// 永远是会保持一致的
	t = time.Now()
	// "Jan 02 2006, 15:04:05" 是golang固定的格式 写错一点 t.Format() 失效
	// M D y, HH:MM:SS
	ctime = t.Format(goParseTimeFormat)

	// 防止撞库攻击 sql 不能使用 + 进行拼接
	//  "INSERT INTO video_info (id,author_id,name,display_ctime,create_time)" +
	//		" VALUES (?,?,?,?,?)" 错误
	var insertSql = `INSERT INTO video_info (id,author_id,name,display_ctime)
		 VALUES (?,?,?,?)`
	if stmtIns,err = dbConn.Prepare(insertSql);err != nil{
		return nil,err
	}
	if _,err = stmtIns.Exec(vid,aid,name,ctime); err != nil {
		return nil,err
	}
	defer stmtIns.Close()
	videoInfo = &model.VideoInfo{
		Id:vid,
		AuthorId:aid,
		Name:name,
		DisplayCtime:ctime,
	}
	return
}

/*
查询video
*/
func GetVideoInfo(vid string) (videoInfo *model.VideoInfo,err error){
	var(
		stmtOut *sql.Stmt
		aid int
		dct string
		name string
	)
	var selectSql = "SELECT author_id, name, display_ctime FROM video_info WHERE id=?"
	if stmtOut,err = dbConn.Prepare(selectSql);err != nil{
		log.Printf("mysql 查询 video_info 表失败 %s",err.Error())
		return nil,err
	}
	if err = stmtOut.QueryRow(vid).Scan(&aid,&name,&dct); err != nil {
		if err != sql.ErrNoRows {
			return nil,err
		}
		if err == sql.ErrNoRows {
			return nil,nil
		}
	}
	defer stmtOut.Close()
	videoInfo = &model.VideoInfo{
		Id:vid,
		AuthorId:aid,
		Name:name,
		DisplayCtime:dct,
	}
	return
}

/*
删除视频
*/
func DeleteVideoInfo(vid string) (err error){
	var(
		stmtDel *sql.Stmt
	)
	var deleteSql = "DELETE FROM video_info WHERE id=?"
	if stmtDel,err = dbConn.Prepare(deleteSql); err != nil {
		log.Printf("mysql 删除 video_info 表失败 %s",err.Error())
		return
	}
	if _,err = stmtDel.Exec(vid); err != nil{
		return
	}
	defer stmtDel.Close()
	return
}

/*
增加评论
@param  vid video id string
@param  aid author id int
@param  content 评论内容 string
@return err
*/
func AddNewComments(vid string, aid int, content string) (err error) {
	var(
		cid string
		stmtIns *sql.Stmt
	)
	if cid,err = common.NewUUID(); err != nil {
		return
	}
	var insertSql = "INSERT INTO comments (id, video_id, author_id, content) values (?, ?, ?, ?)"
	if stmtIns, err = dbConn.Prepare(insertSql); err != nil{
		return
	}
	if _, err = stmtIns.Exec(cid, vid, aid, content); err != nil {
		return
	}
	defer stmtIns.Close()
	return
}

/*
查询评论
*/
func ListComments(vid string, from, to int) (commentList []*model.Comment, err error) {
	var(
		stmtOut *sql.Stmt
		rows *sql.Rows
		comment *model.Comment
	)
	commentList = make([]*model.Comment,0)
	var querySql = ` SELECT comments.id, users.login_name, comments.content 
 					  FROM comments INNER JOIN users 
 					  ON comments.author_id = users.id
					  WHERE comments.video_id = ? 
					  AND comments.time > FROM_UNIXTIME(?) 
					  AND comments.time <= FROM_UNIXTIME(?)`
	if stmtOut,err = dbConn.Prepare(querySql); err != nil{
		log.Printf("mysql 查询 comments 表 失败 %s",err.Error())
		return
	}
	if rows, err = stmtOut.Query(vid, from, to); err != nil {
		return
	}
	for rows.Next() {
		var c_id, login_name, c_content string
		if err = rows.Scan(&c_id, &login_name, &c_content); err != nil {
			return
		}
		comment = &model.Comment{Id: c_id, VideoId: vid, Author: login_name, Content: c_content}
		commentList = append(commentList, comment)
	}
	defer stmtOut.Close()
	return
}