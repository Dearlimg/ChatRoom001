package logic

import (
	"ChatRoom001/dao"
	db "ChatRoom001/dao/mysql/sqlc"
	"ChatRoom001/dao/mysql/tx"
	"ChatRoom001/global"
	"ChatRoom001/model/reply"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type account struct {
}

func (account) CreateAccount(ctx *gin.Context, userID int64, name, avatar, gender, signature string) (*reply.ParamCreateAccount, error) {
	Param := &db.CreateAccountParams{
		ID:        global.GenerateID.GetID(),
		UserID:    userID,
		Name:      name,
		Avatar:    avatar,
		Gender:    db.AccountsGender(gender),
		Signature: signature,
	}

	err := dao.Database.DB.CreateAccountWithTx(ctx, dao.Database.Redis, global.PublicSetting.Rules.AccountNumMax, Param)
	switch {
	case errors.Is(err, tx.ErrAccountOverNum):
	}
}
