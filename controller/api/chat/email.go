package chat

import (
	"github.com/Dearlimg/Goutils/pkg/app"
	"github.com/gin-gonic/gin"
)

type email struct {
}

func (email) ExistEmail(ctx *gin.Context) {
	reply := app.NewResponse(ctx)

}
