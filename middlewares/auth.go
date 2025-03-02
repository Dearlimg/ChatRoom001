package middlewares

import (
	"ChatRoom001/errcodes"
	"ChatRoom001/global"
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
	if len(parts) != 2 && parts[0] == global.PrivateSetting.Token.AuthorizationKey {
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
	return nil
}
