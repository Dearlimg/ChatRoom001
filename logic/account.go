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
	accountInfo, err := GetAccountInfoByID(ctx, userID, accountID)
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
		UserID:     userID,
		Account2ID: ID,
		Account1ID: ID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errcodes.AccountNotFound
		}
		global.Logger.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	return accountInfo, nil
}

func (account) GetAccountsByUserID(ctx *gin.Context, userID int64) (reply.ParamGetAccountByUserID, errcode.Err) {
	accountInfos, err := dao.Database.DB.GetAccountByUserID(ctx, userID)
	if err != nil {
		global.Logger.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return reply.ParamGetAccountByUserID{}, errcode.ErrServer
	}
	result := make([]reply.ParamAccountInfo, len(accountInfos))
	for i, accountInfo := range accountInfos {
		result[i] = reply.ParamAccountInfo{
			ID:     accountInfo.ID,
			Name:   accountInfo.Name,
			Avatar: accountInfo.Avatar,
			Gender: string(accountInfo.Gender),
		}
	}
	return reply.ParamGetAccountByUserID{
		List:  result,
		Total: int64(len(result)),
	}, nil
}

func (account) UpdateAccount(ctx *gin.Context, accountID int64, name, gender, signature string) errcode.Err {
	err := dao.Database.DB.UpdateAccount(ctx, &db.UpdateAccountParams{
		ID:        global.GenerateID.GetID(),
		Name:      name,
		Gender:    db.AccountsGender(gender),
		Signature: signature,
	})
	if err != nil {
		global.Logger.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}

	accessToken, _ := middlewares.GetToken(ctx.Request.Header)
	global.Worker.SendTask(task.UpdateAccount(accessToken, accountID, name, gender, signature))
	return nil
}
