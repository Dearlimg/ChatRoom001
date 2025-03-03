package logic

import (
	"ChatRoom001/model/reply"
	"github.com/Dearlimg/Goutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
)

type user struct {
}

func (user) Register(ctx *gin.Context, emailStr, pwd, code string) (*reply.ParamRegister ,errcode.Err) {
	if err:=CheckEmail

}

//func (U user) Register() error {
//	return nil
//}
