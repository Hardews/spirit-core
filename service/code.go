/**
 * @Author: Hardews
 * @Date: 2023/4/9 20:33
 * @Description:
**/

package service

import (
	"spirit-core/dao"
	"spirit-core/model"
	"spirit-core/my_consts"
	"spirit-core/tool"
)

func CheckCode(addr, code string) error {
	res, err := tool.Get(addr)
	if err != nil {
		return err
	}

	if res != code {
		return my_consts.ErrOfWrongCode
	}

	return dao.AddAddress(model.Address{Address: addr, Status: 0})
}

func SendCode(addr string) error {
	code := RandomStr(5)

	// 存入redis中
	tool.Set(addr, code, my_consts.CodeExp)

	return tool.SendEmail(tool.Theme, addr, code, tool.Code)
}
