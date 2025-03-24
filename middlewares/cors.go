package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
跨域资源共享（CORS）是一种机制，它使用额外的 HTTP 头来告诉浏览器，允许一个 Web 应用在一个域上使用另一个域的资源。
当网页通过脚本访问不属于该网页所在域的资源时，就会发生跨域 HTTP 请求。
*/

// Cors 跨域中间件
//func Cors() gin.HandlerFunc {
//	return func(ctx *gin.Context) {
//		method := ctx.Request.Method               // GET \ POST \ PUT \ DELETE ...
//		origin := ctx.Request.Header.Get("Origin") // 获取请求的 Origin 头部的值，Origin 头部在跨域请求中很重要，会告诉服务器请求的来源域
//		if origin != "" {
//			// 接收客户端发送的 Origin（重要）
//			ctx.Header("Access-Control-Allow-Origin", origin)
//			// 服务器支持的所有跨域请求的方法
//			ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, HEAD, PUT")
//			// 允许跨域设置可以返回其他子段，可以自定义字段
//			ctx.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Origin,content-type,Authorization,Content-Length,X-CSRF-AccessToken,AccessToken,session, token")
//			// 允许浏览器（客户端）可以解析的头部（重要）
//			ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, token")
//			// 允许客户端传递校验信息比如 cookie（重要）
//			ctx.Header("Access-Control-Allow-Credentials", "true")
//		}
//
//		//处理跨域请求时的预检请求，确认服务端是否支持跨域请求，并返回相应的响应。
//		if method == "OPTIONS" {
//			ctx.JSON(http.StatusOK, "options ok")
//			return
//		}
//
//		ctx.Next()
//	}
//}

func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		origin := ctx.Request.Header.Get("Origin")

		// 允许的跨域配置
		if origin != "" {
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin) // 或 "*"
			ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			ctx.Writer.Header().Set("Access-Control-Allow-Headers", "AccountToken, Content-Type, Authorization, AccessToken, X-CSRF-Token")
			ctx.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Content-Type")
			ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

			// 显式允许 WebSocket 升级头
			if ctx.GetHeader("Upgrade") == "websocket" {
				ctx.Writer.Header().Set("Connection", "Upgrade")
				ctx.Writer.Header().Set("Upgrade", "websocket")
			}
		}

		// 处理预检请求
		if method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}

//func Cors() gin.HandlerFunc {
//	return func(ctx *gin.Context) {
//		method := ctx.Request.Method               // GET \ POST \ PUT \ DELETE ...
//		origin := ctx.Request.Header.Get("Origin") // 获取请求的 Origin 头部的值，Origin 头部在跨域请求中很重要，会告诉服务器请求的来源域
//		global.Logger.Info("---------------------------------------------" + origin)
//		if origin != "" {
//			// 接收客户端发送的 Origin（重要）
//			//ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
//			ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
//			// 服务器支持的所有跨域请求的方法
//			ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, HEAD, PUT")
//			// 允许跨域设置可以返回其他子段，可以自定义字段
//			ctx.Writer.Header().Set("Access-Control-Allow-Headers", "AccountToken,Access-Control-Allow-Origin,content-type,AuthorizationKey,Authorization,Content-Length,X-CSRF-AccessToken,AccessToken,session")
//			// 允许浏览器（客户端）可以解析的头部（重要）
//			ctx.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
//			// 允许客户端传递校验信息比如 cookie（重要）
//			ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
//		}
//
//		//处理跨域请求时的预检请求，确认服务端是否支持跨域请求，并返回相应的响应。
//		if method == "OPTIONS" {
//			ctx.JSON(http.StatusOK, "options ok")
//			return
//		}
//
//		ctx.Next()
//	}
//}
