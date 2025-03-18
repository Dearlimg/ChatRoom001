package logic

import (
	"ChatRoom001/global"
	"ChatRoom001/model"
	"github.com/Dearlimg/Goutils/pkg/token"
	"github.com/gin-gonic/gin"
	"time"
)

func newAccountToken(t model.TokenType, id int64) (string, *token.Payload, error) {
	if t != model.AccountToken {
		return "", nil, nil
	}
	duration := global.PrivateSetting.Token.AccountTokenDuration
	data, err := model.NewTokenContent(t, id).Marshal()
	if err != nil {
		return "", nil, err
	}
	result, payload, err := global.TokenMaker.CreateToken(data, duration)
	if err != nil {
		return "", nil, err
	}
	return result, payload, nil
}

// newUserToken
// 成功：返回 token，
func newUserToken(t model.TokenType, id int64, expireTime time.Duration) (string, *token.Payload, error) {
	if t == model.AccountToken {
		return "", nil, nil
	}
	duration := expireTime
	data, err := model.NewTokenContent(t, id).Marshal()
	if err != nil {
		return "", nil, err
	}
	result, payload, err := global.TokenMaker.CreateToken(data, duration)
	if err != nil {
		return "", nil, err
	}
	return result, payload, nil
}

func GetTokenAndPayload(ctx *gin.Context) (string, *token.Payload, error) {
	tokenString := ctx.GetHeader(global.PrivateSetting.Token.AuthorizationType)
	payload, err := global.TokenMaker.VerifyToken(tokenString)
	if err != nil {
		return "", nil, err
	}
	return tokenString, payload, nil
}

// 将 id 从小到大排序返回
func sortID(id1, id2 int64) (_, _ int64) {
	if id1 > id2 {
		return id2, id1
	}
	return id1, id2
}
