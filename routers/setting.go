package routers

import (
	"ChatRoom001/controller/api"
	"ChatRoom001/middlewares"
	"github.com/gin-gonic/gin"
)

type setting struct{}

func (setting) Init(router *gin.RouterGroup) {
	r := router.Group("setting", middlewares.MustAccount())
	{
		r.GET("shows", api.Apis.Setting.GetShows)
		r.GET("pins", api.Apis.Setting.GetPins)
		friendGroup := r.Group("friend")
		{
			friendGroup.GET("list", api.Apis.Setting.GetFriends)
			friendGroup.GET("name", api.Apis.Setting.GetFriendsByName)
		}
		updateGroup := r.Group("update")
		{
			updateGroup.PUT("pin", api.Apis.Setting.UpdateSettingPin)
			updateGroup.PUT("nick_name", api.Apis.Setting.UpdateNickName)
			updateGroup.PUT("show", api.Apis.Setting.UpdateSettingShow)
			updateGroup.PUT("disturb", api.Apis.Setting.UpdateSettingDisturb)
		}
	}
}
