package logic

import (
	"ChatRoom001/global"
	"ChatRoom001/model"
	"fmt"
	"github.com/Dearlimg/Goutils/pkg/token"
	"time"
)

func newAccountToken(t model.TokenType, id int64, expireTime time.Duration) (string, *token.Payload, error) {
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

func newUserToken(t model.TokenType, id int64, expireTime time.Duration) (string, *token.Payload, error) {
	if t == model.AccountToken {
		return "", nil, nil
	}
	duration := expireTime
	fmt.Println("logic common 32 problem in here?")
	data, err := model.NewTokenContent(t, id).Marshal()
	if err != nil {
		return "", nil, err
	}
	fmt.Println("logic common 32 problem in here two?")
	fmt.Println(string(data), duration)
	result, payload, err := global.TokenMaker.CreateToken(data, duration)
	if err != nil {
		return "", nil, err
	}
	return result, payload, nil
}
