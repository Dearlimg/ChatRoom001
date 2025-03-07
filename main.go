package main

import (
	"ChatRoom001/global"
	"ChatRoom001/model/common"
	"ChatRoom001/routers/router"
	"ChatRoom001/setting"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	//"github.com/redis/go-redis/v9"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	//"time"
)

//func testredis() {
//	// 集群模式配置
//	rdb := redis.NewClusterClient(&redis.ClusterOptions{
//		Addrs: []string{
//			"123.249.32.125:6381",
//			"123.249.32.125:6382",
//			"123.249.32.125:6383",
//		},
//		Password:     "1234",          // 集群密码
//		PoolSize:     20,              // 连接池大小
//		ReadOnly:     false,           // 是否使用只读节点
//		DialTimeout:  5 * time.Second, // 连接超时
//		ReadTimeout:  3 * time.Second, // 读取超时
//		WriteTimeout: 3 * time.Second, // 写入超时
//	})
//
//	// 测试连接
//	ctx := context.Background()
//	pong, err := rdb.Ping(ctx).Result()
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println("Redis集群连接成功:", pong)
//	rdb.SAdd(ctx, "chatroom", 1)
//	rdb.SAdd(ctx, "EmailKey", 123)
//
//}

func main() {
	setting.Inits()

	if global.PublicSetting.Server.RunMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("email", common.ValidatorEmail)
	}

	r, ws := router.NewRouter()

	server := &http.Server{
		Addr:           global.PublicSetting.Server.HttpPort,
		Handler:        r,
		MaxHeaderBytes: 1 << 20,
	}
	global.Logger.Info("server start success")
	fmt.Println("AppName:", global.PublicSetting.App.Name, "Version:", global.PublicSetting.App.Version, "Address:", global.PublicSetting.Server.HttpPort, "RunMode:", global.PublicSetting.Server.RunMode)
	errChan := make(chan error, 1)
	defer close(errChan)

	go func() {
		err := server.ListenAndServe()
		//err := r.Run("192.169.1.113:8080")
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer ws.Close()
		if err := ws.Serve(); err != nil {
			errChan <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		global.Logger.Error(err.Error())
	case <-quit:
		global.Logger.Info("Shutdown Server.")
		///创建一个带超时的上下文（给几秒完成还未处理完的请求）
		ctx, cancel := context.WithTimeout(context.Background(), global.PublicSetting.Server.DefaultContextTimeout)
		defer cancel() //延迟取消上下文

		//上下文超时时间内优雅关机（将未处理完的请求处理完再关闭服务），超过超时时间时退出
		if err := server.Shutdown(ctx); err != nil {
			global.Logger.Error("Server forced to Shutdown, err:" + err.Error())
		}
	}
}
