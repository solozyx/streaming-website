package db

import (
	"database/sql"
	"log"
	"strconv"
	"sync"

	"streaming-website/model"
)

// 插入1条session记录
func InsertSession(session_id string,ttl int64,login_name string) (err error){
	var(
		ttlstr string
		stmtIns *sql.Stmt
	)
	ttlstr = strconv.FormatInt(ttl,10)
	var insertSql = `INSERT INTO sessions (session_id,TTL,login_name) VALUES (?,?,?)`
	if stmtIns,err = dbConn.Prepare(insertSql);err != nil{
		return
	}
	defer stmtIns.Close()

	if _,err = stmtIns.Exec(session_id,ttlstr,login_name); err != nil {
		return
	}
	return
}

// 查询1条session记录
func RetrieveSession(session_id string)(session *model.SimpleSession,err error){
	var(
		stmtOut *sql.Stmt
		ttlstr string
		login_name string
		ttlint64 int64
	)
	var selectSql = `SELECT TTL,login_name FROM sessions WHERE session_id = ?`
	if stmtOut,err = dbConn.Prepare(selectSql);err != nil{
		log.Printf("mysql 查询 sessions 表 session_id 字段失败 %s",err.Error())
		return
	}
	if err = stmtOut.QueryRow(session_id).Scan(&ttlstr,&login_name); err != nil {
		if err != sql.ErrNoRows {
			return
		}
		if err == sql.ErrNoRows {
			return
		}
	}
	defer stmtOut.Close()

	// 变量str 进制10 长度64位
	if ttlint64, err = strconv.ParseInt(ttlstr, 10, 64); err == nil {
		session = &model.SimpleSession{
			Username:login_name,
			TTL:ttlint64,
		}
	} else {
		return
	}
	return
}

// 查询所有session记录
// TODO : NOTICE 线程安全 sync.Map 的使用
func RetrieveAllSessions() (sessionList *sync.Map,err error) {
	var(
		stmtOut *sql.Stmt
		rows *sql.Rows
		session_id string
		ttlstr string
		login_name string
		ttlint64 int64
		session *model.SimpleSession
	)

	// make
	sessionList = &sync.Map{}
	var selectSql = "SELECT * FROM sessions"
	if stmtOut, err = dbConn.Prepare(selectSql); err != nil {
		return
	}
	defer stmtOut.Close()

	if rows, err = stmtOut.Query(); err != nil {
		log.Printf("mysql 查询 sessions 表失败 %s", err.Error())
		return
	}

	for rows.Next() {
		if err = rows.Scan(&session_id, &ttlstr, &login_name); err != nil {
			log.Printf("mysql 查询 sessions 表 1条记录转换 失败 %s", err.Error())
			break
		}
		if ttlint64, err = strconv.ParseInt(ttlstr, 10, 64); err == nil{
			session = &model.SimpleSession{
				Username: login_name,
				TTL: ttlint64,
			}
			sessionList.Store(session_id, session)
			log.Printf("session id: %s, ttl: %d", session_id, session.TTL)
		}
	}

	return
}

// 删除1条session记录
func DeleteSession(session_id string) (err error) {
	var(
		stmtOut *sql.Stmt
	)
	var deleteSql = "DELETE FROM sessions WHERE session_id = ?"
	if stmtOut, err = dbConn.Prepare(deleteSql); err != nil {
		log.Printf("mysql 删除 sessions 表 失败%s", err)
		return
	}
	defer stmtOut.Close()

	if _, err = stmtOut.Query(session_id); err != nil {
		return
	}

	return
}