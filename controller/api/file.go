package api

import (
	"github.com/Dearlimg/Goutils/pkg/app"
	"github.com/gin-gonic/gin"
)

type file struct {
}

func (file) PublishFile(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamPublish)
}
