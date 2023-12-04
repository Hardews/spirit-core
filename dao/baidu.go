/**
 * @Author: Hardews
 * @Date: 2023/4/4 23:19
 * @Description:
**/

package dao

import "spirit-core/model"

func GetAccToken() (model.BaiduToken, error) {
	var res model.BaiduToken
	err := db.Model(model.BaiduToken{}).Last(&res).Error

	return res, err
}

func AddAccToken(token model.BaiduToken) error {
	return db.Create(&token).Error
}
