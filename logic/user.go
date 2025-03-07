package logic

import (
	"ChatRoom001/dao"
	db "ChatRoom001/dao/mysql/sqlc"
	"ChatRoom001/errcodes"
	"ChatRoom001/global"
	"ChatRoom001/middlewares"
	"ChatRoom001/model"
	"ChatRoom001/model/reply"
	"fmt"
	"github.com/Dearlimg/Goutils/pkg/app/errcode"
	"github.com/Dearlimg/Goutils/pkg/password"
	"github.com/gin-gonic/gin"
)

type user struct {
}

func (user) Register(ctx *gin.Context, emailStr, pwd, code string) (*reply.ParamRegister, errcode.Err) {
	// 判断邮箱是否已经注册过了
	if err := CheckEmailNotExists(ctx, emailStr); err != nil {
		return nil, err // 已经注册过了
	}
	// 校验验证码
	if !global.EmailMark.CheckCode(emailStr, code) {
		return nil, errcodes.EmailCodeNotValid
	}
	// 密码加密
	hashPassword, err := password.HashPassword(pwd)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	// 将 user 写入数据库并返回 UserInfo
	// 使用MySQL做了一定修改
	err = dao.Database.DB.CreateUser(ctx, &db.CreateUserParams{
		Email:    emailStr,
		Password: hashPassword,
	})
	userInfo, err := dao.Database.DB.GetUserByEmail(ctx, emailStr)

	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	// 添加邮箱到 redis
	err = dao.Database.Redis.AddEmails(ctx, emailStr)
	fmt.Println("logic user 49 : add email success to redis ", emailStr)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	// 创建 token
	accessToken, accessPayload, err := newUserToken(model.UserToken, userInfo.ID, global.PrivateSetting.Token.AccessTokenExpire)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	refreshToken, _, err := newUserToken(model.UserToken, userInfo.ID, global.PrivateSetting.Token.RefreshTokenExpire)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}

	if err = dao.Database.Redis.SaveUserToken(ctx, userInfo.ID, []string{accessToken, refreshToken}); err != nil {
		return nil, errcode.ErrServer.WithDetails(err.Error())
	}

	return &reply.ParamRegister{
		ParamUserInfo: reply.ParamUserInfo{
			ID:       userInfo.ID,
			Email:    userInfo.Email,
			CreateAt: userInfo.CreateAt,
		},
		Token: reply.ParamToken{
			AccessToken:   accessToken,
			AccessPayload: accessPayload,
			RefreshToken:  refreshToken,
		},
	}, nil
}
