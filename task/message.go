package task

import (
	"ChatRoom001/dao"
	"ChatRoom001/global"
	"ChatRoom001/model/chat"
	"ChatRoom001/model/chat/server"
	"ChatRoom001/model/reply"
	"ChatRoom001/pkg/mq/producer"
	"github.com/Dearlimg/Goutils/pkg/utils"
)

func PublishMsg(msg reply.ParamMsgInfoWithRly) func() {
	return func() {
		ctx, cancel := global.DefaultContextWithTimeout()
		defer cancel()
		accountIDs, err := dao.Database.Redis.GetAllAccountsByRelationID(ctx, msg.RelationID)
		if err != nil {
			global.Logger.Error(err.Error())
			return
		}
		for _, accountID := range accountIDs {
			if global.ChatMap.CheckIsOnConnection(accountID) {
				global.ChatMap.Send(accountID, chat.ClientSendMsg, msg)
			} else {
				producer.SendMsgToKafka(accountID, msg)
			}
		}
	}
}

func ReadMsg(accessToken string, readerID int64, msgMap map[int64][]int64, allMsgIDs []int64) func() {
	return func() {
		if len(msgMap) == 0 {
			return
		}
		enToken := utils.EncodeMD5(accessToken)
		for accountID, msgID := range msgMap {
			global.ChatMap.Send(accountID, chat.ClientReadMsg, server.ReadMsg{
				EnToken:  enToken,
				MsgIDs:   msgID,
				ReaderID: readerID,
			})
		}
		global.ChatMap.Send(readerID, chat.ClientReadMsg, server.ReadMsg{
			EnToken:  enToken,
			MsgIDs:   allMsgIDs,
			ReaderID: readerID,
		})
	}
}
