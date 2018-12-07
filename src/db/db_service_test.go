package db

import (
	"testing"
	"fmt"
)

/*
init(dblogin,truncate tables) 每次跑测试表都是空的 初始化状态
-->
run tests
-->
clear data (truncate tables)
避免test和其他包test相互干扰 清理数据-->测试-->清理数据
*/

/*
test初始化
*/
func TestMain(m *testing.M){
	clearTables()
	// 跑所有的test
	m.Run()
	clearTables()
}

func clearTables(){
	dbConn.Exec("TRUNCATE users")
	dbConn.Exec("TRUNCATE video_info")
	dbConn.Exec("TRUNCATE comments")
	dbConn.Exec("TRUNCATE sessions")
}

/*
1条测试路径 规整 子test 调用顺序
*/
func TestUserWorkFlow(t *testing.T){
	t.Run("Add",testAddUser)
	t.Run("Get",testGetUser)
	t.Run("Delete",testDeleteUser)
	t.Run("ReGet",testReGetUser)
}

func testAddUser(t *testing.T) {
	err := AddUserCredential("solozyx","123")
	if err != nil{
		t.Errorf("mysql 插入 user 错误 %s",err)
	}
}

func testGetUser(t *testing.T) {
	pwd,err := GetUserCredential("solozyx")
	if pwd != "123" || err != nil {
		t.Errorf("mysql 查询 user 错误 %s",err)
	}
	fmt.Println("pwd = ",pwd)
}

func testDeleteUser(t *testing.T) {
	err := DeleteUser("solozyx","123")
	if err != nil {
		t.Errorf("mysql 删除 user 错误 %s",err)
	}
}

func testReGetUser(t *testing.T) {
	pwd,err := GetUserCredential("solozyx")
	if err != nil {
		t.Errorf("mysql user 错误 %s",err)
	}
	if pwd != ""{
		t.Errorf("mysql 删除 user 错误 %s",err)
	}
}
