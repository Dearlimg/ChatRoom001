package middlewares

import (
	"ChatRoom001/global"
	"bytes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"time"
)

/*
自定义的 logger 中间件，不使用默认的 logger 中间件
*/

const Body = "body"

// ErrLogMsg 日志数据
func ErrLogMsg(ctx *gin.Context) []zap.Field {
	var body string
	data, ok := ctx.Get("body")
	if ok {
		body = string(data.([]byte))
	}
	path := ctx.Request.URL.Path
	query := ctx.Request.URL.RawQuery
	fields := []zap.Field{
		zap.Int("status", ctx.Writer.Status()),            //记录响应的状态码
		zap.String("method", ctx.Request.Method),          //记录请求方法
		zap.String("path", path),                          //记录请求的路径
		zap.String("query", query),                        //记录请求的原始查询参数
		zap.String("ip", ctx.ClientIP()),                  //记录客户端的 IP 地址
		zap.String("user-agent", ctx.Request.UserAgent()), //记录客户端的 user-agent
		zap.String("body", body),                          //记录请求的主体数据
	}
	return fields
}

// LogBody 读取 body 内容缓存下来，为之后打印日志做准备（读取请求体的内容并将其存储在 Gin 上下文中）
func LogBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyBytes, _ := io.ReadAll(c.Request.Body)
		_ = c.Request.Body.Close()
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		c.Set("body", bodyBytes)
		c.Next()
	}
}

// GinLogger 接收 gin 框架默认的日志，在处理每个请求时记录相关的请求信息到日志中去
func GinLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery
		ctx.Next() //将控制权交给后面的中间件或处理程序执行

		//计算请求处理所耗费的时间
		cost := time.Since(start)
		global.Logger.Info(path,
			zap.Int("status", ctx.Writer.Status()),                                 //记录响应的状态码
			zap.String("method", ctx.Request.Method),                               //记录请求方法
			zap.String("path", path),                                               //记录请求的路径
			zap.String("query", query),                                             //记录请求的原始查询参数
			zap.String("ip", ctx.ClientIP()),                                       //记录客户端的 IP 地址
			zap.String("user-agent", ctx.Request.UserAgent()),                      //记录客户端的 user-agent
			zap.String("errors", ctx.Errors.ByType(gin.ErrorTypePrivate).String()), // ByType 过滤 ctx.Errors 返回的错误列表中指定类型的错误。 gin.ErrorTypePrivate（私有类型的错误）通常由框架或中间件内部使用，不直接向客户端暴露
			zap.Duration("cost", cost),
		)
	}
}
