package api

import (
	"ChatRoom001/errcodes"
	"ChatRoom001/global"
	"ChatRoom001/logic"
	"ChatRoom001/middlewares"
	"ChatRoom001/model"
	"ChatRoom001/model/request"
	"fmt"
	"github.com/Dearlimg/Goutils/pkg/app"
	"github.com/Dearlimg/Goutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
	"time"
)

type message struct{}

func (message) CreateFileMsg(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamCreateFileMsg)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Message.CreateFileMsg(ctx, model.CreateFileMsg{
		AccountID:  content.ID,
		RelationID: params.RelationID,
		RlyMsgID:   params.RlyMsg,
		File:       params.File,
	})
	reply.Reply(err, result)
}

func (message) GetMsgsByRelationIDAndTime(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetMsgsByRelationIDAndTime)
	if err := ctx.ShouldBindBodyWithJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	fmt.Println(params)
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	limit, offset := global.Page.GetPageSizeAndOffset(ctx.Request)
	result, err := logic.Logics.Message.GetMsgsByRelationIDAndTime(ctx, model.GetMsgsByRelationIDAndTime{
		AccountID:  content.ID,
		RelationID: params.RelationID,
		LastTime:   time.Unix(int64(params.LastTime), 0),
		Limit:      limit,
		Offset:     offset,
	})
	reply.Reply(err, result)
}

func (message) OfferMsgsByAccountIDAndTime(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamOfferMsgsByAccountIDAndTime)
	if err := ctx.ShouldBindQuery(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	limit, offset := global.Page.GetPageSizeAndOffset(ctx.Request)
	result, err := logic.Logics.Message.OfferMsgsByAccountIDAndTime(ctx, model.OfferMsgsByAccountIDAndTime{
		AccountID: content.ID,
		LastTime:  time.Unix(int64(params.LastTime), 0),
		Limit:     limit,
		Offset:    offset,
	})
	reply.Reply(err, result)
}

func (message) GetPinMsgsByRelationID(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetPinMsgsByRelationID)
	if err := ctx.ShouldBindQuery(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	limit, offset := global.Page.GetPageSizeAndOffset(ctx.Request)
	result, err := logic.Logics.Message.GetPinMsgsByRelationID(ctx, content.ID, params.RelationID, limit, offset)
	reply.Reply(err, result)
}
