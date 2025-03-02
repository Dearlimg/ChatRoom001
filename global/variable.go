package global

import (
	"ChatRoom001/model/config"
	"ChatRoom001/pkg/emailMark"
	"github.com/Dearlimg/Goutils/pkg/logger"
	"github.com/Dearlimg/Goutils/pkg/token"
)

var (
	Logger         *logger.Log
	EmailMark      *emailMark.EmailMark // 验证码
	PublicSetting  config.PublicConfig
	PrivateSetting config.PrivateConfig //Private 配置
	TokenMaker     token.MakerToken
)
