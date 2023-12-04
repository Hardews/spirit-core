/**
 * @Author: Hardews
 * @Date: 2023/4/5 21:59
 * @Description:
**/

package dao

import (
	"log"
	"time"

	"spirit-core/model"
)

func GetUser(username string) string {
	var user string

	db.Model(&model.User{}).Select("username").Where("username = ?", username).Scan(&user)

	return user
}

func GetPassword(username string) (string, error) {
	err := db.Model(&model.User{}).Select("password").Where("username = ?", username).Scan(&pwd).Error

	return pwd, err
}

func GetRefreshToken(refreshToken string) string {
	var gmtId string

	db.Model(&model.Refresh{}).Select("gmt_id").Where("refresh_token = ? and created_at >= ?", refreshToken, time.Now().Add(-60*60*24*14*time.Second)).Scan(&gmtId)

	return gmtId
}

func GetIdByUsername(username string) string {
	var gmtId string

	db.Model(&model.User{}).Select("gmt_id").Where("username = ?", username).Scan(&gmtId)

	return gmtId
}

func GetAdmIdentity(gmtId string) model.User {
	var info model.User
	db.Model(model.User{}).Where("gmt_id = ?", gmtId).First(&info)
	return info
}

func AddRefreshToken(refreshToken, gmtId string) error {
	return db.Create(&model.Refresh{
		RefreshToken: refreshToken,
		GmtId:        gmtId,
	}).Error
}

func AddUser(user model.User) error {
	return db.Create(&user).Error
}

func DelRefreshToken(rt, id string) {
	err := db.Where("refresh_token = ? AND gmt_id = ?", rt, id).Delete(&model.Refresh{})
	if err != nil {
		log.Println("del refresh token failed,err:", err)
	}
}

func UpdatePassword(password, gmtId string) error {
	return db.Model(&model.User{}).Where("gmt_id = ?", gmtId).
		Update("password", password).Error
}
