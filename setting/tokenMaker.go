package setting

import (
	"ChatRoom001/global"
	"github.com/Dearlimg/Goutils/pkg/token"
)

type tokenMaker struct {
}

func (tokenMaker) Init() {
	var err error
	global.TokenMaker, err = token.NewPasetoMaker([]byte(global.PrivateSetting.Token.Key))
	if err != nil {
		panic(err)
	}
}
