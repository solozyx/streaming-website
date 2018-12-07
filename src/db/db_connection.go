package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

var(
	dbConn *sql.DB
	err error
)

/*
整个程序中当db包被加载的时候init()会第一个被执行
*/
func init(){
	// golang内置sql接口 常用的数据库都有实现
	dbConn,err = sql.Open(
		"mysql",
		"root:root@tcp(localhost:3306)/video_info?charset=utf8")
	// defer dbConn.Close()
	if err != nil {
		goto ERR
	}
	return
ERR:
	fmt.Println(err.Error())
}