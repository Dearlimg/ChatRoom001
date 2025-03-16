package task

import (
	"ChatRoom001/global"
	"ChatRoom001/model/chat"
	"ChatRoom001/model/chat/server"
	"github.com/Dearlimg/Goutils/pkg/utils"
)

/*
有关 account 的 任务
*/

//func UpdateEmail(accountToken string, userID int64, email string) func() {
//	return func() {
//		ctx, cancel := global.DefaultContextWithTimeout()
//		defer cancel()
//		accountIDs, err := dao.Database.DB.GetAcountIDsByUserID(ctx, userID)
//		if err != nil {
//			global.Logger.Error(err.Error())
//			return
//		}
//		global.Chat
//	}
//}

func UpdateAccount(accessToken string, accountID int64, name, gender, signature string) func() {
	return func() {
		global.ChatMap.Send(accountID, chat.ServerUpdateAccount, server.UpdateAccount{
			EnToken:   utils.EncodeMD5(accessToken),
			Name:      name,
			Gender:    gender,
			Signature: signature,
		})
	}
}
