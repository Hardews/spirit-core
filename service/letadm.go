/**
 * @Author: Hardews
 * @Date: 2023/3/16 10:33
 * @Description:
**/

package service

import (
	"log"
	"spirit-core/dao"
	"spirit-core/model"
	"strconv"
)

func UpdateLetterStatus(id, status int) error {
	return dao.UpdateLetterStatus(id, status)
}

func GetLetterStatus(id string) (int, error) {
	return dao.GetLetterStatus(id)
}

func GetLetterHomepage() ([]model.LetterOverview, error) {
	var ls []model.LetterOverview

	letters, err := dao.GetLetterOverview()
	if err != nil {
		return nil, err
	}

	for _, letter := range letters {
		ls = append(ls, model.LetterOverview{
			Id:         letter.ID,
			Reply:      letter.Reply,
			Time:       letter.CreatedAt.Unix(),
			Status:     letter.Status,
			Predicted:  letter.Predicted,
			Conclusion: letter.Conclusion,
			Msg:        letter.Msg,
		})
	}

	return ls, err
}

func GetLetterInfo(id string) (lr model.LetterReturn, err error) {
	var letterReply string

	letter, err := dao.GetLetterInfo(id)
	if err != nil {
		return
	}

	// 如果已经被回复
	if letter.Reply {
		letterReply, err = dao.GetReplyInfo(id)
		if err != nil {
			return
		}
	}

	lr = model.LetterReturn{
		Id:         letter.ID,
		Email:      letter.Address,
		Text:       letter.Context,
		Reply:      letterReply,
		Time:       letter.CreatedAt.Unix(),
		Audit:      letter.Reply,
		Status:     letter.Status,
		Predicted:  letter.Predicted,
		Conclusion: letter.Conclusion,
		Msg:        letter.Msg,
	}

	return
}

func ReplyLetter(letterId int, admId, content string) (err error) {
	letterInfo, err := dao.GetLetterInfo(strconv.Itoa(letterId))
	if err != nil {
		log.Println("get letter info failed,err:", err)
	}

	storage := model.LetterReply{
		LetterContent: letterInfo.Context,
		LetterId:      letterId,
		Context:       content,
		GmtId:         admId,
		Audit:         0,
	}

	return dao.AddLetterReply(storage)
}
