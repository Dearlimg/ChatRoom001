package setting

import (
	"ChatRoom001/global"
	"github.com/Dearlimg/Goutils/pkg/generateID/snowflake"
	"time"
)

type generateID struct {
}

func (generateID) Init() {
	var err error
	global.GenerateID, err = snowflake.Init(time.Now(), global.PublicSetting.App.MachineID)
	if err != nil {
		panic(err)
	}
}
