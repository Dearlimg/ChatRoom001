package logic

import (
	"ChatRoom001/dao"
	"ChatRoom001/errcodes"
	"ChatRoom001/global"
	"ChatRoom001/middlewares"
	"ChatRoom001/model/reply"
	"ChatRoom001/pkg/emailMark"
	"errors"
	"math/rand"
	"strings"

	"github.com/Dearlimg/Goutils/pkg/app/errcode"

	"github.com/gin-gonic/gin"
)

type email struct {
}

// ExistEmail 是否存在 email
func (email) ExistEmail(ctx *gin.Context, emailStr string) (*reply.ParamExistEmail, errcode.Err) {
	// 先在 redis 缓存中查找
	//dao.Database.Redis.TestEmailRedis(ctx, emailStr)
	ok, err := dao.Database.Redis.ExistEmail(ctx, emailStr)
	if err == nil {
		return &reply.ParamExistEmail{Exist: ok}, nil
	}
	global.Logger.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
	// 如果在 redis 中没找到，再到 MySQL 数据库中查找
	ok, err = dao.Database.DB.ExistEmail(ctx, emailStr)
	if err != nil {
		global.Logger.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	return &reply.ParamExistEmail{Exist: ok}, nil
}

// CheckEmailNotExists 判断邮箱是否存在
// 先去缓存中查询，如果不在缓存中，再去数据库中查询，如果存在，将邮箱写入缓存中，返回邮箱已注册的错误。（旁路缓存中的读策略）
func CheckEmailNotExists(ctx *gin.Context, emailStr string) errcode.Err {
	result, err := email{}.ExistEmail(ctx, emailStr)
	if err != nil {
		return err
	}
	if result.Exist {
		return errcodes.EmailExists
	}
	exist, myErr := dao.Database.DB.ExistEmail(ctx, emailStr)
	if myErr != nil {
		global.Logger.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	if exist {
		myErr = dao.Database.Redis.AddEmails(ctx, emailStr)
		if myErr != nil {
			global.Logger.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
			return errcode.ErrServer
		}
		return errcodes.EmailExists
	}
	return nil
}

func RandomNumberString(n int) string {
	var sb strings.Builder
	digits := "0123456789" // 数字字符集
	k := len(digits)

	for i := 0; i < n; i++ {
		c := digits[rand.Intn(k)] // 随机选择一个数字字符
		sb.WriteByte(c)
	}
	return sb.String()
}

// SendMark 发送验证码(邮件)
func (email) SendMark(emailStr string) errcode.Err {
	// 判断发送邮件的频率
	if global.EmailMark.CheckUserExist(emailStr) {
		return errcodes.EmailSendMany
	}
	// 异步发送邮件(使用工作池)
	global.Worker.SendTask(func() {
		//code := utils.RandomString(global.PublicSetting.Rules.CodeLength)
		code := RandomNumberString(global.PublicSetting.Rules.CodeLength)
		if err := global.EmailMark.SendMark(emailStr, code); err != nil && !errors.Is(err, emailMark.ErrSendTooMany) {
			//fmt.Println("logic/email.go")
			global.Logger.Error(err.Error())
		}
	})
	return nil
}
