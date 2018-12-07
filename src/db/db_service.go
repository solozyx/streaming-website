package db

import (
	"database/sql"
	"log"
)

func AddUserCredential(loginName string,pwd string) (err error) {
	var(
		stmtIns *sql.Stmt
	)
	if stmtIns,err = dbConn.Prepare("INSERT INTO users (login_name,pwd) VALUES (?,?)");err != nil{
		return
	}
	stmtIns.Exec(loginName,pwd)
	// 提高性能 这里 不写 defer stmtIns.Close()
	stmtIns.Close()
	return
}

func GetUserCredential(loginName string) (pwd string,err error) {
	var(
		stmtOut *sql.Stmt
	)
	if stmtOut,err = dbConn.Prepare("SELECT pwd FROM users WHERE login_name = ?");err != nil{
		log.Printf("mysql 查询 user 表 pwd 字段失败 %s",err.Error())
		return "",err
	}
	// query 1 row data
	stmtOut.QueryRow(loginName).Scan(&pwd)
	stmtOut.Close()
	return pwd,nil
}

func DeleteUser(loginName string,pwd string) (err error){
	var(
		stmtDel *sql.Stmt
	)
	if stmtDel,err = dbConn.Prepare("DELETE FROM users WHERE login_name=? and pwd=?"); err != nil {
		log.Printf("mysql 删除 user 表 loginName 字段失败 %s",err.Error())
		return
	}
	stmtDel.Exec(loginName,pwd)
	stmtDel.Close()
	return
}