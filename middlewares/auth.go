package middlewares

import (
	"ChatRoom001/dao"
	"ChatRoom001/errcodes"
	"ChatRoom001/global"
	"ChatRoom001/model"
	"fmt"
	"github.com/Dearlimg/Goutils/pkg/app"
	"github.com/Dearlimg/Goutils/pkg/app/errcode"
	"github.com/Dearlimg/Goutils/pkg/token"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

/*
用户验证（paseTo 生成 Token）
*/

// GetToken 从当前请求头获取 token
func GetToken(header http.Header) (string, errcode.Err) {
	authorizationHeader := header.Get(global.PrivateSetting.Token.AuthorizationKey)
	if len(authorizationHeader) == 0 {
		return "", errcodes.AuthNotExist
	}
	parts := strings.SplitN(authorizationHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == global.PrivateSetting.Token.AuthorizationType) {
		return "", errcodes.AuthenticationFailed
	}
	return parts[1], nil
}

// ParseToken 解析 header 中的 token。返回 payload，token，err
func ParseToken(token string) (*token.Payload, string, errcode.Err) {
	payload, err := global.TokenMaker.VerifyToken(token)
	if err != nil {
		if err.Error() == "超时错误" {
			return nil, "", errcodes.AuthOverTime
		}
		return nil, "", errcodes.AuthenticationFailed
	}
	return payload, token, nil
}

// PasetoAuth 鉴权中间件，用于解析并写入 Token
func PasetoAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		accessToken, err := GetToken(c.Request.Header)
		if err != nil {
			//fmt.Println("PasetoAuth1:")
			c.Next()
			return
		}
		payload, _, err := ParseToken(accessToken)
		if err != nil {
			//fmt.Println("PasetoAuth2:")
			c.Next()
			return
		}
		content := &model.Content{}
		if err := content.Unmarshal(payload.Content); err != nil {
			c.Next()
			return
		}
		c.Set(global.PrivateSetting.Token.AuthorizationKey, content)
		fmt.Println("PasetoAuth:", content, payload)
		c.Next()
	}
}

func MustUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		reply := app.NewResponse(c)
		val, ok := c.Get(global.PrivateSetting.Token.AuthorizationKey)
		if !ok {
			reply.Reply(errcodes.AuthNotExist)
			c.Abort()
			return
		}
		data := val.(*model.Content)

		if data.TokenType != model.UserToken {
			reply.Reply(errcodes.AuthenticationFailed)
			c.Abort()
			return
		}
		//fmt.Println(data)
		ok, err := dao.Database.DB.ExistsUserByID(c, data.ID)
		if err != nil {
			global.Logger.Error(err.Error(), ErrLogMsg(c)...)
			reply.Reply(errcode.ErrServer)
			c.Abort()
			return
		}
		if !ok {
			reply.Reply(errcodes.UserNotFound)
			c.Abort()
			return
		}
		c.Next()
	}
}

func MustAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		reply := app.NewResponse(c)
		val, ok := c.Get(global.PrivateSetting.Token.AuthorizationKey)
		if !ok {
			reply.Reply(errcodes.AuthNotExist)
			c.Abort()
			return
		}
		data := val.(*model.Content)
		fmt.Println("MustAccount", data)
		if data.TokenType != model.AccountToken {
			reply.Reply(errcodes.AuthenticationFailed)
			c.Abort()
			return
		}
	}
}

func GetTokenContent(ctx *gin.Context) (*model.Content, bool) {
	value, ok := ctx.Get(global.PrivateSetting.Token.AuthorizationKey)
	if !ok {
		return nil, false
	}
	return value.(*model.Content), true
}
