package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"streaming-website/conf"
)

var(
	dbConn *sql.DB
	err error
)

// 整个程序中当db包被加载的时候init()会第一个被执行
func init(){
	// golang内置sql接口 常用的数据库都有实现
	dbConn,err = sql.Open("mysql",conf.DBConnUrl)
	// defer dbConn.Close()
	if err != nil {
		goto ERR
	}
	return
ERR:
	fmt.Println(err.Error())
}