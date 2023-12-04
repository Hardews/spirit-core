/**
 * @Author: Hardews
 * @Date: 2023/4/5 21:58
 * @Description:
**/

package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	GmtId    string
	Username string
	Nickname string
	Password string
	Address  string
}

type Refresh struct {
	gorm.Model
	GmtId        string
	RefreshToken string
}

type UserInfoReturn struct {
	Identity string
	Nickname string
	Address  string
}
