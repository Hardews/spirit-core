/**
 * @Author: Hardews
 * @Date: 2023/3/16 21:16
 * @Description:
**/

package dao

import (
	"spirit-core/model"
)

func GetLetterStatus(id string) (int, error) {
	var status int
	err := db.Set("gorm:query_option", "FOR UPDATE").Model(&model.Letter{}).Select("status").Where("id = ?", id).Scan(&status).Error
	if err != nil {
		return 0, err
	}

	return status, err
}

func GetLetterOverview() ([]model.Letter, error) {
	var ls []model.Letter

	err := db.Model(&model.Letter{}).Order("id DESC").Scan(&ls).Error

	return ls, err
}

// GetLetterInfo 获取信件详细内容
func GetLetterInfo(id string) (letter model.Letter, err error) {
	err = db.Model(model.Letter{}).Where("id = ?", id).First(&letter).Error
	return
}

// GetReplyInfo 如果回信了获取回信的文本
func GetReplyInfo(letterId string) (string, error) {
	var reply string

	err := db.Model(&model.LetterReply{}).Select("context").Where("letter_id = ?", letterId).First(&reply).Error

	return reply, err
}

func AddLetterReply(reply model.LetterReply) error {
	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()

	err := tx.Create(&reply).Error
	if err != nil {
		return err
	}

	err = tx.Model(&model.Letter{}).Where("id = ?", reply.LetterId).UpdateColumns(map[string]interface{}{"status": 2, "gmt_id": reply.GmtId}).Error
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func UpdateLetterStatus(id, status int) error {
	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()

	err := tx.Model(&model.Letter{}).Where("id = ?", id).Update("status", status).Error
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}
