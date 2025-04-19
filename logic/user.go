package logic

import (
	"ChatRoom001/dao"
	db "ChatRoom001/dao/mysql/sqlc"
	"ChatRoom001/errcodes"
	"ChatRoom001/global"
	"ChatRoom001/middlewares"
	"ChatRoom001/model"
	"ChatRoom001/model/reply"
	"ChatRoom001/task"
	"database/sql"
	"github.com/Dearlimg/Goutils/pkg/app/errcode"
	"github.com/Dearlimg/Goutils/pkg/password"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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

func (user) Login(ctx *gin.Context, emailStr, pwd string) (*reply.ParamLogin, errcode.Err) {
	userInfo, err := dao.Database.DB.GetUserByEmail(ctx, emailStr)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
	}
	if err := password.CheckPassword(pwd, userInfo.Password); err != nil {
		return nil, errcodes.PasswordNotValid
	}
	//token
	accessToken, accessPayload, err := newUserToken(model.UserToken, userInfo.ID, global.PrivateSetting.Token.AccessTokenExpire)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	reflashToken, _, err := newUserToken(model.UserToken, userInfo.ID, global.PrivateSetting.Token.RefreshTokenExpire)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	if err = dao.Database.Redis.SaveUserToken(ctx, userInfo.ID, []string{accessToken, reflashToken}); err != nil {
		return nil, errcode.ErrServer.WithDetails(err.Error())
	}
	return &reply.ParamLogin{
		ParamUserInfo: reply.ParamUserInfo{
			ID:       userInfo.ID,
			Email:    userInfo.Email,
			CreateAt: userInfo.CreateAt,
		},
		Token: reply.ParamToken{
			AccessToken:   accessToken,
			AccessPayload: accessPayload,
			RefreshToken:  reflashToken,
		},
	}, nil

}

func (user) Logout(ctx *gin.Context) errcode.Err {
	Token, payload, err := GetTokenAndPayload(ctx)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcodes.AuthenticationFailed
	}
	content := &model.Content{}
	_ = content.Unmarshal(payload.Content)
	//判断是否再redis
	if ok := dao.Database.Redis.CheckUserTokenValid(ctx, content.ID, Token); !ok {
		return errcodes.UserNotFound
	}
	if err := dao.Database.Redis.DeleteAllTokenByUser(ctx, content.ID); err != nil {
		return errcode.ErrServer.WithDetails(err.Error())
	}
	return nil
}

func getUserInfoByID(ctx *gin.Context, userID int64) (*db.User, errcode.Err) {
	userInfo, err := dao.Database.DB.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errcodes.UserNotFound
		}
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	return userInfo, nil
}

func (user) UpdateUserPassword(ctx *gin.Context, userID int64, code, newPwd string) errcode.Err {
	userInfo, myerr := getUserInfoByID(ctx, userID)
	if myerr != nil {
		global.Logger.Error(myerr.Error(), middlewares.ErrLogMsg(ctx)...)
		return myerr
	}
	// 校验验证码
	if !global.EmailMark.CheckCode(userInfo.Email, code) {
		return errcodes.EmailCodeNotValid
	}
	// 密码加密
	hashPassword, err := password.HashPassword(newPwd)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	// 更新密码
	if err = dao.Database.DB.UpdateUser(ctx, &db.UpdateUserParams{
		Email:    userInfo.Email,
		Password: hashPassword,
		ID:       userInfo.ID,
	}); err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}

	// 清除用户的 token
	if err := dao.Database.Redis.DeleteAllTokenByUser(ctx, userID); err != nil {
		return errcode.ErrServer.WithDetails(err.Error())
	}
	return nil
}

func (user) UpdateUserEmail(ctx *gin.Context, userID int64, email, code string) errcode.Err {
	userInfo, myerr := getUserInfoByID(ctx, userID)
	if myerr != nil {
		global.Logger.Error(myerr.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	if userInfo.Email == email {
		return errcodes.EmailSame
	}
	if err := CheckEmailNotExists(ctx, email); err != nil {
		return err
	}
	if !global.EmailMark.CheckCode(userInfo.Email, code) {
		return errcodes.EmailCodeNotValid
	}
	if err := dao.Database.DB.UpdateUser(ctx, &db.UpdateUserParams{
		Email:    email,
		Password: userInfo.Password,
		ID:       userInfo.ID,
	}); err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	if err := dao.Database.Redis.UpdateEmail(ctx, userInfo.Email, email); err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	accessToken, _ := middlewares.GetToken(ctx.Request.Header)
	global.Worker.SendTask(task.UpdateEmail(accessToken, userID, email))
	return nil
}

func (user) DeleteUser(ctx *gin.Context, UserID int64) errcode.Err {
	userInfo, myerr := getUserInfoByID(ctx, UserID)
	if myerr != nil {
		global.Logger.Error(myerr.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	accountNum, err := dao.Database.DB.CountAccountByUserID(ctx, UserID)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	if accountNum > 0 {
		return errcodes.UserHasAccount
	}
	if err := dao.Database.DB.DeleteUser(ctx, UserID); err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	if err := dao.Database.Redis.DeleteEmail(ctx, userInfo.Email); err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	Token, payload, err := GetTokenAndPayload(ctx)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcodes.AuthenticationFailed
	}
	content := &model.Content{}
	_ = content.Unmarshal(payload.Content)
	if ok := dao.Database.Redis.CheckUserTokenValid(ctx, content.ID, Token); !ok {
		return errcodes.UserNotFound
	}
	//先将token从redis中清除
	if err := dao.Database.Redis.DeleteAllTokenByUser(ctx, content.ID); err != nil {
		return errcode.ErrServer.WithDetails(err.Error())
	}
	return nil
}
