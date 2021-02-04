package model

// 用户信息
type UserDetails struct {
	// 用户标示
	UserId int
	// 用户名 唯一
	Username string
	// 用户密码
	Password string
	// 用户具有的权限
	Authorities []string
}
