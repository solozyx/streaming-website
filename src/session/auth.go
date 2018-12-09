package session

import (
	"net/http"
)

var(
	// X- 自定义header field
	HEADER_FIELD_SESSION_ID = "X-Session-Id"
	HEADER_FIELD_USER_NAME  = "X-User-Name"
)

/*
校验用户session是否合法
*/
func ValidateUserSession(req *http.Request) (isUserSessionValid bool){
	var(
		session_id string
		login_name string
		ok bool
	)
	session_id = req.Header.Get(HEADER_FIELD_SESSION_ID)
	if len(session_id) == 0 {
		return false
	}
	if login_name,ok = IsSessionExpired(session_id); ok {
		// session过期
		return false
	}
	req.Header.Add(HEADER_FIELD_USER_NAME,login_name)
	return true
}

/*
校验用户是否合法
*/
func ValidateUser(resp http.ResponseWriter,req *http.Request) (isUserValid bool){
	var(
		login_name string
	)
	if login_name = req.Header.Get(HEADER_FIELD_USER_NAME); len(login_name) == 0 {
		return false
	}
	return true
}