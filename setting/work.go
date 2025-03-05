package setting

import (
	"ChatRoom001/global"
	"github.com/Dearlimg/Goutils/pkg/goroutine/work"
)

type worker struct {
}

func (worker) Init() {
	global.Worker = work.Init(work.Config{
		TaskChanCapacity:   global.PublicSetting.Worker.TaskChanCapacity,
		WorkerChanCapacity: global.PublicSetting.Worker.WorkerChanCapacity,
		WorkerNum:          global.PublicSetting.Worker.WorkerNum,
	})
}
