/**
 * @Author: Hardews
 * @Date: 2023/4/5 21:58
 * @Description:
**/

package service

import (
	"log"
	"spirit-core/dao"
	"spirit-core/model"
	"spirit-core/my_consts"
	"spirit-core/tool"
)

func GetToken(user model.User) (token, refreshToken string, err error) {
	// 判断用户是否存在
	if !IsUserHas(user.Username) {
		err = my_consts.ErrOfNoThisUser
		return
	}

	gmtId := dao.GetIdByUsername(user.Username)
	if gmtId == "" {
		err = my_consts.ErrOfNoThisUser
		return
	}

	if IsPasswordRight(user.Username, user.Password) {
		// 颁发token
		return tool.TokenGenerate(gmtId)
	} else {
		err = my_consts.ErrOfWrongPassword
		return
	}
}

func RefreshToken(rt string) (token, refreshToken string, err error) {
	return tool.RefreshToken(rt)
}

// IsUserHas 判断是否有这个用户
func IsUserHas(username string) bool {
	return dao.GetUser(username) == username
}

func IsPasswordRight(username, password string) bool {
	// 判断密码是否正确
	pwd, err := dao.GetPassword(username)
	if err != nil {
		log.Println("get password failed,err:", err)
		return false
	}

	// 默认密码
	if pwd == "" {
		password = password + my_consts.DefaultSalt
		pwd, _ = tool.Encryption(username[:8] + my_consts.DefaultSalt)
	}

	return tool.CheckPassword(pwd, password)
}

func ChangePassword(user model.User, oldPassword, gmtId string) (bool, error) {
	var err error

	if !IsUserHas(user.Username) {
		return true, my_consts.ErrOfNoThisUser
	}

	IsPasswordRight(user.Username, oldPassword)
	if !IsPasswordRight(user.Username, oldPassword) {
		return true, my_consts.ErrOfWrongPassword
	}

	user.Password, err = tool.Encryption(user.Password)
	if err = dao.UpdatePassword(user.Password, gmtId); err != nil {
		return false, err
	}

	return true, nil
}

func GetIdentity(gmtId string) model.UserInfoReturn {
	info := dao.GetAdmIdentity(gmtId)
	identity := "admin"
	if dao.IsPermission(gmtId) {
		identity = "super-adm"
	}

	return model.UserInfoReturn{
		Identity: identity,
		Nickname: info.Nickname,
		Address:  info.Address,
	}
}

func AddUser(gmtId string, user model.User) error {
	if !dao.IsPermission(gmtId) {
		return my_consts.ErrOfForbidden
	}

	user.Password, _ = tool.Encryption(user.Password)
	user.GmtId = tool.GenerateGMTId(user.Address, user.Username)

	return dao.AddUser(user)
}
