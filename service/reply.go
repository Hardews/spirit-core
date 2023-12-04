/**
 * @Author: Hardews
 * @Date: 2023/3/25 14:36
 * @Description:
**/

package service

import (
	"spirit-core/dao"
	"spirit-core/model"
	"spirit-core/my_consts"
	"spirit-core/tool"
)

func ChangeLetterReplyStatus(replyId int, admId, content string, isPass bool) error {
	// 先查看是否有权限
	if !dao.IsPermission(admId) {
		return my_consts.ErrOfForbidden
	}

	// 是否通过审核 是->发送邮件给用户 否->驳回并存储信息
	if !isPass {
		// 这里修改letter的status 为 0，并修改 adm 的表
		return dao.AddReplyAdvice(replyId, admId, content)
	}

	// 通过审核,给用户发送邮件
	// 更新数据库
	err := dao.UpdateLetterReplyStatus(replyId, admId)
	if err != nil {
		return err
	}

	// 从数据库中读取回信的信息
	addr, reply, err := dao.GetReplyUserInfo(replyId)
	if err != nil {
		return err
	}

	return tool.SendEmail(tool.Theme, addr, reply, tool.Reply)
}

func GetMyReplyInfo(replyId, gmtId string) (model.LetterReply, error) {
	return dao.GetMyReplyLetterInfo(replyId, gmtId)
}

func GetReplyInfo(replyId, gmtId string) (model.LetterReply, error) {
	if !dao.IsPermission(gmtId) {
		return model.LetterReply{}, my_consts.ErrOfForbidden
	}

	return dao.GetReplyLetterInfo(replyId)
}

func GetReplyLetterHomepage(admId string) ([]model.LetterReplyOverview, error) {
	if !dao.IsPermission(admId) {
		return nil, my_consts.ErrOfForbidden
	}

	var ls []model.LetterReplyOverview

	letters, err := dao.GetReplyOverview()
	if err != nil {
		return nil, err
	}

	for _, letter := range letters {
		ls = append(ls, model.LetterReplyOverview{
			Id:     letter.ID,
			Status: letter.Audit,
			Time:   letter.CreatedAt.Unix(),
		})
	}

	return ls, err
}

func GetMyReplyLetter(admId string) ([]model.LetterReplyOverview, error) {
	var ls []model.LetterReplyOverview

	letters, err := dao.GetMyReplyOverview(admId)
	if err != nil {
		return nil, err
	}

	for _, letter := range letters {
		ls = append(ls, model.LetterReplyOverview{
			Id:     letter.ID,
			Status: letter.Audit,
			Time:   letter.CreatedAt.Unix(),
		})
	}

	return ls, err
}

func GetReplySum() int64 {
	return dao.GetReplySum()
}
