package setting

import (
	"ChatRoom001/global"
	"github.com/Dearlimg/Goutils/pkg/app"
)

type page struct{}

func (page) Init() {
	global.Page = app.InitPage(global.PublicSetting.Page.DefaultPageSize, global.PublicSetting.Page.MaxPageSize, global.PublicSetting.Page.PageKey, global.PublicSetting.Page.PageSizeKey)
}
