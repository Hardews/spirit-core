/**
 * @Author: Hardews
 * @Date: 2023/3/16 10:07
 * @Description:
**/

package model

import (
	"gorm.io/gorm"
	"time"
)

// Letter 信件相关结构体
type Letter struct {
	gorm.Model
	Context    string  // 信件内容
	Address    string  // 回复的邮件地址
	GmtID      string  // 分配回复的管理员id
	Reply      bool    // 是否已经被回复
	Status     int     // 状态 0 未处理 1 有人在处理 2 已经被处理了
	Predicted  float64 // 预测值，取值区间[-3,3]
	Conclusion string  // 是否合规
	Msg        string  // 为什么不合规
}

// LetterInfo 信件其他信息存储的相关结构体
type LetterInfo struct {
	gorm.Model
	Text         string // 信件内容
	Confidence   float64
	NegativeProb float64
	PositiveProb float64
	Sentiment    int
	LogId        int64 // 情感分析时的logId
}

// Future 未来信相关结构体
type Future struct {
	gorm.Model
	Name        string    `json:"name,omitempty"`    // 名字
	Text        string    `json:"text,omitempty"`    // 内容
	Theme       string    `json:"theme,omitempty"`   // 主题
	Address     string    `json:"address,omitempty"` // 发送的地址
	SendTime    time.Time `json:"send_time"`         // 发送时间,格式 YY/MM/DD 如 2023/4/14 23:39
	IsPublic    bool      `json:"is_public"`         // 是否公开 0 不公开 1 公开
	IsSend      bool      `json:"is_send"`           // 是否发送 1 未发送 0 已发送
	ShouldAudit bool      `json:"-"`                 // 是否需要审核 0 需要 1 不需要
}

type Address struct {
	gorm.Model
	LetterId int
	Address  string
	Status   int
}

type LetterFutureForm struct {
	Content  string `json:"content,omitempty"`
	Name     string `json:"send_name,omitempty"`
	Theme    string `json:"theme,omitempty"`
	Address  string `json:"address,omitempty"`
	Time     string `json:"send_time,omitempty"`
	IsPublic bool   `json:"is_public,omitempty"`
}

type LetterReturn struct {
	Id         uint    `json:"id"`
	Email      string  `json:"email"`
	Text       string  `json:"text"`
	Reply      string  `json:"reply"`
	Time       int64   `json:"time"`
	Audit      bool    `json:"audit"`
	Status     int     `json:"status"`
	Predicted  float64 `json:"predicted"`
	Conclusion string  `json:"conclusion"`
	Msg        string  `json:"msg"`
}

type LetterOverview struct {
	Id         uint    `json:"id"`
	Reply      bool    `json:"reply"`
	Time       int64   `json:"time"`
	Status     int     `json:"status"`
	Predicted  float64 `json:"predicted"`
	Conclusion string  `json:"conclusion"`
	Msg        string  `json:"msg"`
}
