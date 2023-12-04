/**
 * @Author: Hardews
 * @Date: 2023/3/25 14:36
 * @Description:
**/

package dao

import (
	"spirit-core/model"

	"gorm.io/gorm"
)

func IsPermission(admId string) bool {
	var check string
	err := db.Model(&model.Authority{}).Select("gmt_id").Where("gmt_id = ?", admId).First(&check).Error

	return check == admId || err != gorm.ErrRecordNotFound
}

func GetMyReplyLetterInfo(replyId, gmtId string) (model.LetterReply, error) {
	var reply model.LetterReply

	err := db.Model(&model.LetterReply{}).Where("id = ? AND gmt_id = ?", replyId, gmtId).Scan(&reply).Error

	return reply, err
}

func GetReplyLetterInfo(replyId string) (model.LetterReply, error) {
	var reply model.LetterReply

	err := db.Model(&model.LetterReply{}).Where("id = ?", replyId).Scan(&reply).Error

	return reply, err
}

func GetMyReplyOverview(adm string) ([]model.LetterReply, error) {
	var ls []model.LetterReply

	err := db.Model(&model.LetterReply{}).Where("gmt_id = ?", adm).Order("id DESC").Scan(&ls).Error

	return ls, err
}

func GetReplyOverview() ([]model.LetterReply, error) {
	var ls []model.LetterReply

	err := db.Model(&model.LetterReply{}).Order("id DESC").Scan(&ls).Error

	return ls, err
}

// GetReplyUserInfo 获取回信时用户的信息
func GetReplyUserInfo(replyId int) (string, string, error) {
	var addr string
	var replyInfo struct {
		LetterId int
		Context  string
	}

	err := db.Model(&model.LetterReply{}).Select([]string{"letter_id", "context"}).Where("id = ?", replyId).Scan(&replyInfo).Error
	if err != nil || replyInfo.LetterId == 0 {
		return "", "", err
	}

	err = db.Model(&model.Letter{}).Select("address").Where("id = ?", replyInfo.LetterId).Scan(&addr).Error

	return addr, replyInfo.Context, err
}

func GetReplySum() int64 {
	var sum int64
	db.Model(&model.LetterReply{}).Count(&sum)
	return sum
}

func AddReplyAdvice(replyId int, admId, content string) error {
	var letterId int
	db.Model(&model.LetterReply{}).Select("letter_id").Where("id = ?", replyId).Scan(&letterId)
	db.Model(&model.Letter{}).Where("id = ?", letterId).Update("status", 0)
	return db.Model(&model.LetterReply{}).Where("id = ?", replyId).
		UpdateColumns(map[string]interface{}{"advice": content, "Audit": 2}).Error
}

func UpdateLetterReplyStatus(replyId int, admId string) error {
	var letterId int
	db.Model(&model.LetterReply{}).Select("letter_id").Where("id = ?", replyId).First(&letterId)
	if letterId == 0 {
		return gorm.ErrRecordNotFound
	}

	var audit int
	db.Model(&model.LetterReply{}).Select("audit").Where("id = ?", replyId).Scan(&audit)
	if audit == 1 {
		return nil
	}

	tx := db.Begin()
	defer tx.Rollback()

	// 通过replyId修改参数
	err := tx.Model(&model.LetterReply{}).Where("id = ?", replyId).Update("audit", 1).Error
	if err != nil {
		return err
	}

	// 通过letterId修改相关数值
	err = tx.Model(&model.Letter{}).Where("id = ?", letterId).Update("reply", true).Error
	if err != nil {
		return err
	}

	tx.Commit()

	// 异步插入操作记录
	go db.Create(&model.SuperOperate{
		ReplyId: replyId,
		GmtId:   admId,
	})

	return nil
}
