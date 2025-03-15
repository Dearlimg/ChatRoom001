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
	fmt.Println(token, payload, err)
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
