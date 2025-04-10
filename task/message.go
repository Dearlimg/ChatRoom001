package task

import (
	"ChatRoom001/dao"
	"ChatRoom001/global"
	"ChatRoom001/model/chat"
	"ChatRoom001/model/chat/server"
	"ChatRoom001/model/reply"
	"ChatRoom001/pkg/mq/producer"
	"fmt"
	"github.com/Dearlimg/Goutils/pkg/utils"
)

// PublishMsg 推送消息事件和执行拓展内容
// 参数：消息和回复消息
func PublishMsg(msg reply.ParamMsgInfoWithRly) func() {
	return func() {
		ctx, cancel := global.DefaultContextWithTimeout()
		defer cancel()
		accountIDs, err := dao.Database.Redis.GetAllAccountsByRelationID(ctx, msg.RelationID)
		//fmt.Println("PublishMsg accountIDs ", accountIDs)
		if err != nil {
			global.Logger.Error(err.Error())
			return
		}
		count := 0
		fmt.Println("\u001B[32mPublishMsg time sd aw] :  \u001B[0m\n", accountIDs, msg.RelationID)
		for _, accountID := range accountIDs {
			// 用户如果在线，直接将消息发送过去
			if global.ChatMap.CheckIsOnConnection(accountID) {
				fmt.Println("\u001B[32mPublishMsg time :  \u001B[0m\n", count)
				count++
				global.ChatMap.Send(accountID, chat.ClientSendMsg, msg)
			} else { // 用户处于离线状态，将消息发送至 MQ 中
				fmt.Printf("\033[32m[Put message to mq] \033[0m\n")
				producer.SendMsgToKafka(accountID, msg)
			}
		}
		fmt.Println("\u001B[32mPublishMsg time  finished:  \u001B[0m\n", count)
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

func UpdateMsgState(accessToken string, relationID, msgID int64, msgType server.MsgType, state bool) func() {
	return func() {
		ctx, cancel := global.DefaultContextWithTimeout()
		defer cancel()
		accountIDs, err := dao.Database.Redis.GetAllAccountsByRelationID(ctx, relationID)
		if err != nil {
			global.Logger.Error(err.Error())
			return
		}
		global.ChatMap.SendMany(accountIDs, chat.ServerUpdateMsgState, server.UpdateMsgState{
			EnToken: utils.EncodeMD5(accessToken),
			MsgType: msgType,
			MsgID:   msgID,
			State:   state,
		})
	}
}
