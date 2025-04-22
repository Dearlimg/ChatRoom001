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
	if err := ctx.ShouldBindQuery(params); err != nil {
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
	//fmt.Println("GetMsgsByRelationIDAndTime test NULL data", result)
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
	fmt.Println("GetPinMsgsByRelati dwa wd 1 onID:", content.ID, params, params.RelationID, limit, offset)
	result, err := logic.Logics.Message.GetPinMsgsByRelationID(ctx, content.ID, params.RelationID, limit, offset)
	reply.Reply(err, result)
}

func (message) GetRlyMsgsInfoByMsgID(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetRlyMsgsInfoByMsgID)
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
	result, err := logic.Logics.Message.GetRlyMsgsInfoByMsgID(ctx, content.ID, params.RelationID, params.MsgID, limit, offset)
	reply.Reply(err, result)
}

func (message) GetMsgsByContent(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetMsgsByContent)
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
	result, err := logic.Logics.Message.GetMsgsByContent(ctx, content.ID, params.RelationID, params.Content, limit, offset)
	reply.ReplyList(err, result.Total, result)
}

func (message) GetTopMsgByRelationID(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetTopMsgByRelationID)
	if err := ctx.ShouldBindQuery(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Message.GetTopMsgByRelationID(ctx, content.ID, params.RelationID)
	reply.Reply(err, result)
}

func (message) UpdateMsgPin(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamUpdateMsgPin)

	if err := ctx.ShouldBindJSON(params); err != nil {
		//fmt.Println("UpdateMsgPin  ", params, err)
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}

	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}

	err := logic.Logics.Message.UpdateMsgPin(ctx, content.ID, params)
	reply.Reply(err)
}

func (message) UpdateMsgTop(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamUpdateMsgTop)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Message.UpdateMsgTop(ctx, content.ID, params)
	reply.Reply(err)
}

func (message) RevokeMsg(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamRevokeMsg)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Message.RevokeMsg(ctx, content.ID, params.ID)
	reply.Reply(err)
}
