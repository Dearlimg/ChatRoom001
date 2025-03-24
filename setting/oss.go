package setting

import (
	"ChatRoom001/global"
	"github.com/Dearlimg/Goutils/pkg/upload/obs/ali_cloud"
)

type oss struct{}

func (oss) Init() {
	global.OSS = ali_cloud.Init(ali_cloud.Config{
		Location:         global.PrivateSetting.HuaWeiOBS.Location,
		BucketName:       global.PrivateSetting.HuaWeiOBS.BucketName,
		BucketUrl:        global.PrivateSetting.HuaWeiOBS.BucketUrl,
		Endpoint:         global.PrivateSetting.HuaWeiOBS.Endpoint,
		BasePath:         global.PrivateSetting.HuaWeiOBS.BasePath,
		AvatarType:       global.PrivateSetting.HuaWeiOBS.AvatarType,
		AccountAvatarUrl: global.PrivateSetting.HuaWeiOBS.AccountAvatarUrl,
		GroupAvatarUrl:   global.PrivateSetting.HuaWeiOBS.GroupAvatarUrl,
	})
}
