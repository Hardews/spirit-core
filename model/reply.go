/**
 * @Author: Hardews
 * @Date: 2023/3/25 14:37
 * @Description:
**/

package model

import "gorm.io/gorm"

// LetterReply 信件回复相关结构体
type LetterReply struct {
	gorm.Model    `json:"-"`
	LetterContent string `json:"letter_content"`
	LetterId      int    `json:"letter_id,omitempty"` // 信件的id
	Context       string `json:"context,omitempty"`   // 信件回复内容
	Advice        string `json:"advice"`              // 信件回复未通过
	GmtId         string `json:"-"`                   // 回复信件的管理员id
	Audit         int    `json:"audit"`               // 该回复信件是否通过审核  0->未审核 1->审核通过 2->审核未通过
}

type SuperOperate struct {
	gorm.Model
	ReplyId int
	GmtId   string
}

type Authority struct {
	gorm.Model
	GmtId string
}

type LetterReplyOverview struct {
	Id     uint  `json:"id,omitempty"`
	Status int   `json:"status"`
	Time   int64 `json:"time"`
}
