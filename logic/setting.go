package logic

import (
	"ChatRoom001/dao"
	db "ChatRoom001/dao/mysql/sqlc"
	"ChatRoom001/errcodes"
	"ChatRoom001/global"
	"ChatRoom001/middlewares"
	"ChatRoom001/model"
	"ChatRoom001/model/chat/server"
	"ChatRoom001/model/reply"
	"ChatRoom001/task"
	"context"
	"database/sql"
	"fmt"
	"github.com/Dearlimg/Goutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type setting struct{}

func ExistsSetting(ctx context.Context, accountID, relationID int64) (bool, errcode.Err) {
	ok, err := dao.Database.DB.ExistsSetting(ctx, &db.ExistsSettingParams{
		AccountID:  accountID,
		RelationID: relationID,
	})
	if err != nil {
		global.Logger.Error(err.Error())
		return false, errcode.ErrServer
	}
	return ok, nil
}

func (setting) GetFriends(ctx *gin.Context, accountID int64) (*reply.ParamGetFriends, errcode.Err) {
	data, err := dao.Database.DB.GetFriendSettingsOrderByName(ctx, &db.GetFriendSettingsOrderByNameParams{
		AccountID:   accountID,
		AccountID_2: accountID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	result := make([]*model.SettingFriend, 0, len(data))
	for _, v := range data {
		result = append(result, &model.SettingFriend{
			SettingInfo: model.SettingInfo{
				RelationID:   v.RelationID,
				RelationType: "friend",
				NickName:     v.NickName,
				IsNotDisturb: v.IsNotDisturb,
				IsPin:        v.IsPin,
				IsShow:       v.IsShow,
				PinTime:      v.PinTime,
				LastShow:     v.LastShow,
			},
			FriendInfo: &model.SettingFriendInfo{
				AccountID: v.AccountID,
				Name:      v.AccountName,
				Avatar:    v.AccountAvatar,
			},
		})
	}
	return &reply.ParamGetFriends{
		List:  result,
		Total: int64(len(result)),
	}, nil
}

func (setting) GetFriendsByName(ctx *gin.Context, accountID int64, name string, limit, offset int32) (*reply.ParamGetFriendsByName, errcode.Err) {
	data, err := dao.Database.DB.GetFriendSettingsByName(ctx, &db.GetFriendSettingsByNameParams{
		AccountID:   accountID,
		AccountID_2: accountID,
		Limit:       limit,
		Offset:      offset,
		CONCAT:      name,
		CONCAT_2:    name,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	if len(data) == 0 {
		return &reply.ParamGetFriendsByName{List: []*model.SettingFriend{}}, nil
	}
	result := make([]*model.SettingFriend, 0, len(data))
	for _, v := range data {
		result = append(result, &model.SettingFriend{
			SettingInfo: model.SettingInfo{
				RelationID:   v.RelationID,
				RelationType: string(db.RelationsRelationTypeFriend),
				NickName:     v.NickName,
				IsNotDisturb: v.IsNotDisturb,
				IsPin:        v.IsPin,
				IsShow:       v.IsShow,
				PinTime:      v.PinTime,
				LastShow:     v.LastShow,
			},
			FriendInfo: &model.SettingFriendInfo{
				AccountID: v.AccountID,
				Name:      v.AccountName,
				Avatar:    v.AccountAvatar,
			},
		})
	}
	return &reply.ParamGetFriendsByName{
		List:  result,
		Total: data[0].Total.(int64),
	}, nil
}

func (setting) GetShows(ctx *gin.Context, accountID int64) (*reply.ParamGetShows, errcode.Err) {
	friendData, err := dao.Database.DB.GetFriendShowSettingsOrderByShowTime(ctx, &db.GetFriendShowSettingsOrderByShowTimeParams{
		AccountID:   accountID,
		AccountID_2: accountID,
	})
	//friendData, err := dao.Database.DB.GetFriendShowSettingsOrderByShowTime(ctx, accountID)

	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}

	groupData, err := dao.Database.DB.GetGroupShowSettingsOrderByShowTime(ctx, &db.GetGroupShowSettingsOrderByShowTimeParams{
		AccountID:   accountID,
		AccountID_2: accountID,
	})

	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	//fmt.Println("GetShows 3 ", err)
	result := make([]*model.Setting, 0, len(friendData)+len(groupData))
	for i, j := 0, 0; i < len(friendData) || j < len(groupData); {
		if i < len(friendData) && (j >= len(groupData) || friendData[i].Lastshow.After(groupData[j].LastShow)) {
			v := friendData[i]
			fmt.Println("GetShow data", v)
			//msgInfo, myErr := dao.Database.DB.GetLastMessageByRelation(ctx, v.RelationID)
			//if myErr != nil {
			//	global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
			//	return nil, errcode.ErrServer
			//}
			friendInfo := &model.SettingFriendInfo{
				AccountID: v.Accountid,
				Name:      v.Accountname,
				Avatar:    v.Accountavatar,
				Create_at: v.Accountcreateat,
			}
			result = append(result, &model.Setting{
				SettingInfo: model.SettingInfo{
					RelationID:   v.Relationid,
					RelationType: "friend",
					NickName:     v.Nickname,
					IsNotDisturb: v.Isnotdisturb,
					IsPin:        v.Ispin,
					IsShow:       v.Isshow,
					PinTime:      v.Pintime,
					LastShow:     v.Lastshow,

					//Msg_id:      msgInfo.ID,
					//Msg_type:    string(msgInfo.MsgType),
					//Msg_content: msgInfo.MsgContent,
					//Create_at:   msgInfo.CreateAt,
				},
				FriendInfo: friendInfo,
			})
			i++
		} else {
			v := groupData[j]
			//groupType := strings.Split(v., ",")
			//msgInfo, myErr := dao.Database.DB.GetLastMessageByRelation(ctx, v.RelationID)
			//if myErr != nil {
			//	global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
			//	return nil, errcode.ErrServer
			//}
			//fmt.Println("GetShows 423323 ", err, msgInfo)
			//if msgInfo.ID == 0 || msgInfo.MsgType == "" || msgInfo.MsgContent == "" {
			//	msgInfo.MsgType = "text"
			//	msgInfo.ID = 0
			//	msgInfo.MsgContent = "nil"
			//}
			groupInfo := &model.SettingGroupInfo{
				RelationID:  v.RelationID,
				Name:        v.GroupName.String,
				Description: v.GroupDescription.String,
				Avatar:      v.GroupAvatar.String,
			}
			result = append(result, &model.Setting{
				SettingInfo: model.SettingInfo{
					RelationID:   v.RelationID,
					RelationType: "group",
					NickName:     v.NickName,
					IsNotDisturb: v.IsNotDisturb,
					IsPin:        v.IsPin,
					IsShow:       v.IsShow,
					PinTime:      v.PinTime,
					LastShow:     v.LastShow,

					//Msg_id:      msgInfo.ID,
					//Msg_type:    string(msgInfo.MsgType),
					//Msg_content: msgInfo.MsgContent,
					//Create_at:   msgInfo.CreateAt,
				},
				GroupInfo: groupInfo,
			})
			j++
		}
	}
	fmt.Println("GetShows 4 ")
	return &reply.ParamGetShows{
		List:  result,
		Total: int64(len(result)),
	}, nil
}

func (setting) GetPins(ctx *gin.Context, accountID int64) (*reply.ParamGetPins, errcode.Err) {
	friendData, err := dao.Database.DB.GetFriendPinSettingsOrderByPinTime(ctx, &db.GetFriendPinSettingsOrderByPinTimeParams{
		AccountID:   accountID,
		AccountID_2: accountID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	groupData, err := dao.Database.DB.GetGroupPinSettingsOrderByPinTime(ctx, &db.GetGroupPinSettingsOrderByPinTimeParams{
		AccountID:   accountID,
		AccountID_2: accountID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	result := make([]*model.SettingPin, 0, len(groupData)+len(friendData))
	for i, j := 0, 0; i < len(friendData) || j < len(groupData); {
		if i < len(friendData) && (j >= len(groupData) || friendData[i].PinTime.Before(groupData[j].PinTime)) {
			v := friendData[i]
			msgInfo, myErr := dao.Database.DB.GetLastMessageByRelation(ctx, v.RelationID)
			if myErr != nil {
				global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
				return nil, errcode.ErrServer
			}
			friendInfo := &model.SettingFriendInfo{
				AccountID:      accountID,
				Name:           v.AccountName,
				Avatar:         v.AccountAvatar,
				Is_Pin:         v.IsPin,
				Is_Show:        v.IsShow,
				Is_Not_Disturb: v.IsNotDisturb,

				Msg_id:      msgInfo.ID,
				Msg_type:    string(msgInfo.MsgType),
				Msg_content: msgInfo.MsgContent,
				Create_at:   msgInfo.CreateAt,
			}
			result = append(result, &model.SettingPin{
				SettingPinInfo: model.SettingPinInfo{
					RelationID:   v.RelationID,
					RelationType: "friend",
					NickName:     v.NickName,
					PinTime:      v.PinTime,
				},
				FriendInfo: friendInfo,
			})
			i++
		} else {
			v := groupData[j]
			msgInfo, myErr := dao.Database.DB.GetLastMessageByRelation(ctx, v.RelationID)
			if myErr != nil {
				global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
				return nil, errcode.ErrServer
			}
			groupInfo := &model.SettingGroupInfo{
				RelationID:     v.RelationID,
				Name:           v.GroupName_2.String,
				Description:    v.GroupDescription_2.String,
				Avatar:         v.GroupAvatar.String,
				Is_Pin:         v.IsPin,
				Is_Show:        v.IsShow,
				Is_Not_Disturb: v.IsNotDisturb,

				Msg_id:      msgInfo.ID,
				Msg_type:    string(msgInfo.MsgType),
				Msg_content: msgInfo.MsgContent,
				Create_at:   msgInfo.CreateAt,
			}
			result = append(result, &model.SettingPin{
				SettingPinInfo: model.SettingPinInfo{
					RelationID:   v.RelationID,
					RelationType: "group",
					NickName:     v.NickName,
					PinTime:      v.PinTime,
				},
				GroupInfo: groupInfo,
			})
			j++
		}
	}
	return &reply.ParamGetPins{
		List:  result,
		Total: int64(len(result)),
	}, nil
}

func (setting) UpdateNickName(ctx *gin.Context, accountID, relationID int64, nickName string) errcode.Err {
	settingInfo, err := dao.Database.DB.GetSettingByID(ctx, &db.GetSettingByIDParams{
		AccountID:  accountID,
		RelationID: relationID,
	})
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return errcodes.RelationNotExists
	case errors.Is(err, nil):
		if settingInfo.NickName == nickName {
			return nil
		}
		if err := dao.Database.DB.UpdateSettingNickName(ctx, &db.UpdateSettingNickNameParams{
			NickName:   nickName,
			AccountID:  accountID,
			RelationID: relationID,
		}); err != nil {
			global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
			return errcode.ErrServer
		}
		// 向自己推送更改昵称的通知
		accessToken, _ := middlewares.GetToken(ctx.Request.Header)
		global.Worker.SendTask(task.UpdateNickName(accessToken, accountID, relationID, nickName))
		return nil
	default:
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
}

func (setting) UpdatePin(ctx *gin.Context, accountID, relationID int64, isPin bool) errcode.Err {
	settingInfo, err := dao.Database.DB.GetSettingByID(ctx, &db.GetSettingByIDParams{
		AccountID:  accountID,
		RelationID: relationID,
	})
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return errcodes.RelationNotExists
	case errors.Is(err, nil):
		if settingInfo.IsPin == isPin {
			return nil
		}
		if err := dao.Database.DB.UpdateSettingPin(ctx, &db.UpdateSettingPinParams{
			IsPin:      isPin,
			AccountID:  accountID,
			RelationID: relationID,
		}); err != nil {
			global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
			return errcode.ErrServer
		}
		accessToken, _ := middlewares.GetToken(ctx.Request.Header)
		global.Worker.SendTask(task.UpdateSettingState(accessToken, server.SettingPin, accountID, relationID, isPin))
		return nil
	default:
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
}

func (setting) UpdateShow(ctx *gin.Context, accountID, relationID int64, isShow bool) errcode.Err {
	settingInfo, err := dao.Database.DB.GetSettingByID(ctx, &db.GetSettingByIDParams{
		AccountID:  accountID,
		RelationID: relationID,
	})
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return errcodes.RelationNotExists
	case errors.Is(err, nil):
		if settingInfo.IsShow == true {
			return nil
		}
		err = dao.Database.DB.UpdateSettingShow(ctx, &db.UpdateSettingShowParams{
			IsShow:     isShow,
			AccountID:  accountID,
			RelationID: relationID,
		})
		if err != nil {
			global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
			return errcode.ErrServer
		}
		accessToken, _ := middlewares.GetToken(ctx.Request.Header)
		global.Worker.SendTask(task.UpdateSettingState(accessToken, server.SettingShow, accountID, relationID, isShow))
		return nil
	default:
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
}

func (setting) UpdateDisaturb(ctx *gin.Context, accountID, relationID int64, isDisat bool) errcode.Err {
	settingInfo, err := dao.Database.DB.GetSettingByID(ctx, &db.GetSettingByIDParams{
		AccountID:  accountID,
		RelationID: relationID,
	})
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return errcodes.RelationNotExists
	case errors.Is(err, nil):
		if settingInfo.IsNotDisturb == isDisat {
			return nil
		}
		err = dao.Database.DB.UpdateSettingDisturb(ctx, &db.UpdateSettingDisturbParams{
			IsNotDisturb: isDisat,
			AccountID:    accountID,
			RelationID:   relationID,
		})
		if err != nil {
			global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
			return errcode.ErrServer
		}
		accessToken, _ := middlewares.GetToken(ctx.Request.Header)
		global.Worker.SendTask(task.UpdateSettingState(accessToken, server.SettingNotDisturb, accountID, relationID, isDisat))
		return nil
	default:
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
}
