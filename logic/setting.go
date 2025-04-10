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
	result := make([]*model.Setting, 0, len(groupData))
	for i, j := 0, 0; i < len(friendData) || j < len(groupData); {
		if i < len(friendData) && (j >= len(groupData) || friendData[i].LastShow.After(groupData[j].LastShow)) {
			v := friendData[i]
			friendInfo := &model.SettingFriendInfo{
				AccountID: v.AccountID,
				Name:      v.AccountName,
				Avatar:    v.AccountAvatar,
			}
			result = append(result, &model.Setting{
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
				FriendInfo: friendInfo,
			})
			i++
		} else {
			v := groupData[j]
			//groupType := strings.Split(v., ",")
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
				},
				GroupInfo: groupInfo,
			})
			j++
		}
	}
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
			friendInfo := &model.SettingFriendInfo{
				AccountID: accountID,
				Name:      v.AccountName,
				Avatar:    v.AccountAvatar,
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
			groupInfo := &model.SettingGroupInfo{
				RelationID:  v.RelationID,
				Name:        v.GroupName_2.String,
				Description: v.GroupDescription_2.String,
				Avatar:      v.GroupAvatar.String,
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
