package global

import (
	"ChatRoom001/model/config"
	"ChatRoom001/pkg/emailMark"
	"github.com/Dearlimg/Goutils/pkg/generateID/snowflake"
	"github.com/Dearlimg/Goutils/pkg/goroutine/work"
	"github.com/Dearlimg/Goutils/pkg/logger"
	"github.com/Dearlimg/Goutils/pkg/token"
)

var (
	Logger         *logger.Log
	EmailMark      *emailMark.EmailMark // 验证码
	PublicSetting  config.PublicConfig
	PrivateSetting config.PrivateConfig //Private 配置
	TokenMaker     token.MakerToken
	Worker         *work.Worker
	GenerateID     *snowflake.Snowflake //snowflake 雪花算法生成的 ID
)
