package db

import (
	"database/sql"
	"log"
	"model"
)

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
		log.Printf("mysql 删除 user 表 loginName 字段失败 %s",err.Error())
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
@return video obj
*/
func AddNewVideo(aid int,name string) (videoInfo *model.VideoInfo,err error){

	return
}