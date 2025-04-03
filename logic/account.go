package logic

import (
	"ChatRoom001/dao"
	db "ChatRoom001/dao/mysql/sqlc"
	"ChatRoom001/dao/mysql/tx"
	"ChatRoom001/errcodes"
	"ChatRoom001/global"
	"ChatRoom001/middlewares"
	"ChatRoom001/model"
	"ChatRoom001/model/common"
	"ChatRoom001/model/reply"
	"ChatRoom001/task"
	"database/sql"
	"fmt"
	"github.com/Dearlimg/Goutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type account struct {
}

func (account) CreateAccount(ctx *gin.Context, userID int64, name, avatar, gender, signature string) (*reply.ParamCreateAccount, errcode.Err) {
	Param := &db.CreateAccountParams{
		ID:        global.GenerateID.GetID(),
		UserID:    userID,
		Name:      name,
		Avatar:    avatar,
		Gender:    db.AccountsGender(gender),
		Signature: signature,
	}

	err := dao.Database.DB.CreateAccountWithTx(ctx, dao.Database.Redis, int64(global.PublicSetting.Rules.AccountNumMax), Param)
	switch {
	case errors.Is(err, tx.ErrAccountOverNum):
		return nil, errcodes.AccountNumExcessive
	case errors.Is(err, tx.ErrAccountNameExist):
		return nil, errcodes.AccountNameExists
	case err == nil:
	default:
		global.Logger.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	token, payload, err := newAccountToken(model.AccountToken, Param.ID)
	if err != nil {
		global.Logger.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	return &reply.ParamCreateAccount{
		ParamAccountInfo: reply.ParamAccountInfo{
			ID:     Param.ID,
			Name:   name,
			Avatar: avatar,
			Gender: gender,
		},
		ParamGetAccountToken: reply.ParamGetAccountToken{
			AccountToken: common.Token{
				Token:    token,
				ExpireAt: payload.ExpiredAt,
			},
		},
	}, nil
}

func (account) GetAccountToken(ctx *gin.Context, userID, accountID int64) (*reply.ParamGetAccountToken, errcode.Err) {
	accountInfo, err := GetAccountInfoByID(ctx, accountID, accountID)
	if err != nil {
		return nil, err
	}
	if accountInfo.UserID != userID {
		return nil, errcodes.AuthPermissionsInsufficient
	}
	token, payload, err1 := newAccountToken(model.AccountToken, accountID)
	if err1 != nil {
		global.Logger.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	return &reply.ParamGetAccountToken{AccountToken: common.Token{
		Token:    token,
		ExpireAt: payload.ExpiredAt,
	}}, nil
}

func GetAccountInfoByID(ctx *gin.Context, userID, accountID int64) (*db.GetAccountByIDRow, errcode.Err) {
	var ID sql.NullInt64
	ID.Int64 = accountID

	accountInfo, err := dao.Database.DB.GetAccountByID(ctx, &db.GetAccountByIDParams{
		ID:         userID,
		Account2ID: ID,
		Account1ID: ID,
	})
	//fmt.Println("GetAccountInfoByID ", userID, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errcodes.AccountNotFound
		}
		global.Logger.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	return accountInfo, nil
}

//func GetAccountInfoByID(ctx *gin.Context, accountID, selfID int64) (*db.GetAccountByIDRow, errcode.Err) {
//		var ID sql.NullInt64
//		ID.Int64 = accountID
//
//	accountInfo, err := dao.Database.DB.GetAccountByID(ctx, &db.GetAccountByIDParams{
//				ID:         selfID,
//				Account2ID: ID,
//				Account1ID: ID,
//	})
//	if err != nil {
//		if errors.Is(err, pgx.ErrNoRows) {
//			return nil, errcodes.AccountNotFound
//		}
//		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
//		return nil, errcode.ErrServer
//	}
//	return accountInfo, nil
//}

func (account) GetAccountsByUserID(ctx *gin.Context, userID int64) (reply.ParamGetAccountByUserID, errcode.Err) {
	accountInfos, err := dao.Database.DB.GetAccountByUserID(ctx, userID)
	if err != nil {
		global.Logger.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return reply.ParamGetAccountByUserID{}, errcode.ErrServer
	}
	result := make([]reply.ParamAccountInfo, len(accountInfos))
	for i, accountInfo := range accountInfos {
		result[i] = reply.ParamAccountInfo{
			ID:        accountInfo.ID,
			Name:      accountInfo.Name,
			Avatar:    accountInfo.Avatar,
			Gender:    string(accountInfo.Gender),
			Signature: accountInfo.Signature,
		}
	}
	return reply.ParamGetAccountByUserID{
		List:  result,
		Total: int64(len(result)),
	}, nil
}

func (account) UpdateAccount(ctx *gin.Context, accountID int64, name, gender, signature string) errcode.Err {
	err := dao.Database.DB.UpdateAccount(ctx, &db.UpdateAccountParams{
		Name:      name,
		Gender:    db.AccountsGender(gender),
		Signature: signature,
		ID:        accountID,
	})
	if err != nil {
		global.Logger.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		//fmt.Println("is here", ctx, accountID, name, gender, signature, err)
		return errcode.ErrServer
	}

	accessToken, _ := middlewares.GetToken(ctx.Request.Header)
	global.Worker.SendTask(task.UpdateAccount(accessToken, accountID, name, gender, signature))
	return nil
}

func (account) GetAccountsByName(ctx *gin.Context, accountID int64, name string, limit, offset int32) (reply.ParamGetAccountsByName, errcode.Err) {
	var ID sql.NullInt64
	ID.Int64 = accountID
	accounts, err := dao.Database.DB.GetAccountsByName(ctx, &db.GetAccountsByNameParams{
		Limit:      limit,
		Offset:     offset,
		CONCAT:     name,
		Account1ID: ID,
		Account2ID: ID,
	})
	fmt.Println("account153", accountID, name, limit, offset)
	if err != nil {
		global.Logger.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return reply.ParamGetAccountsByName{}, errcode.ErrServer
	}
	result := make([]*reply.ParamFriendInfo, len(accounts))
	for i, info := range accounts {
		result[i] = &reply.ParamFriendInfo{
			ParamAccountInfo: reply.ParamAccountInfo{
				ID:     info.ID,
				Name:   info.Name,
				Avatar: info.Avatar,
				Gender: string(info.Gender),
			},
			RelationID: info.RelationID.Int64,
		}
	}
	return reply.ParamGetAccountsByName{
		List:  result,
		Total: int64(len(result)),
	}, nil
}

func (account) GetAccountByID(ctx *gin.Context, accountID, selfID int64) (*reply.ParamGetAccountByID, errcode.Err) {
	fmt.Println("GetAccountByID ACCOUNTID:", accountID, "SELFID", selfID)
	info, err := GetAccountInfoByID(ctx, selfID, accountID)
	if err != nil {
		return nil, err
	}
	return &reply.ParamGetAccountByID{
		Info: reply.ParamAccountInfo{
			ID:     info.ID,
			Name:   info.Name,
			Avatar: info.Avatar,
			Gender: string(info.Gender),
		},
		Signature:  info.Signature,
		CreateAt:   info.CreateAt,
		RelationID: info.RelationID.Int64,
	}, nil
}

func (account) DeleteAccount(ctx *gin.Context, selfID int64, accountID int64) errcode.Err {
	accountInfo, myerr := GetAccountInfoByID(ctx, selfID, accountID)
	if myerr != nil {
		return myerr
	}
	if accountInfo.UserID != selfID {
		return errcodes.AuthPermissionsInsufficient
	}
	err := dao.Database.DB.DeleteAccountWithTx(ctx, dao.Database.Redis, accountID)
	switch {
	case errors.Is(err, tx.ErrAccountGroupLeader):
		return errcodes.AccountGroupLeader
	case errors.Is(err, nil):
		return nil
	default:
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
}
