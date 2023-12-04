/**
 * @Author: Hardews
 * @Date: 2023/4/5 21:19
 * @Description:
**/

package dao

import (
	"gorm.io/gorm"
	"spirit-core/model"
	"time"
)

var LetterSum int64

func IsPermAddr(addr string) bool {
	var check string
	err := db.Model(&model.Address{}).Select("address").Where("address = ? AND status = 0", addr).Last(&check).Error

	return check == addr || err != gorm.ErrRecordNotFound
}

func AddLetter(letter model.Letter, info model.LetterInfo) error {
	tx := db.Begin()
	defer tx.Rollback()

	err := tx.Create(&letter).Error
	if err != nil {
		return err
	}

	err = tx.Create(&info).Error
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func AddAddress(addr model.Address) error {
	return db.Create(&addr).Error
}

func AddFutureLetter(fl model.Future) error {
	return db.Create(&fl).Error
}

func UpdateAddress(addr string) {
	db.Model(&model.Address{}).Where("address = ?", addr).Last(&model.Address{}).Update("status", 1)
}

func UpdateFutureStatus(id uint) error {
	return db.Model(&model.Future{}).Where("id = ?", id).Update("is_send", true).Error
}

func GetTodayFutureEmail() ([]model.Future, error) {
	var fls []model.Future
	err := db.Model(&model.Future{}).Where("send_time BETWEEN ? AND ?", time.Now(), time.Now().Add(time.Hour*24)).Scan(&fls).Error

	return fls, err
}

func GetFutureEmail() []model.Future {
	var fls []model.Future
	db.Model(&model.Future{}).Where("is_public = true").Order("id DESC").Scan(&fls)
	return fls
}

func GetFutureEmailById(id string) model.Future {
	var fls model.Future
	db.Model(&model.Future{}).Where("id = ? AND is_public = true", id).Scan(&fls)
	return fls
}

func GetLetterSum() int64 {
	var sum int64
	db.Model(&model.Letter{}).Count(&sum)
	return sum
}

func GetLetterFutureSum() int64 {
	var sum int64
	db.Model(&model.Future{}).Count(&sum)
	return sum
}

func GetFutureLetterPublicSum() int64 {
	var sum int64
	db.Model(&model.Future{}).Where("is_public = true").Count(&sum)
	LetterSum = sum
	return sum
}

func GetFutureIds() []int {
	var ids []int
	db.Model(&model.Future{}).Where("is_public = true").Select("id").Scan(&ids)
	return ids
}
