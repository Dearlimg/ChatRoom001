package retry

import (
	"fmt"
	"time"
)

type Try struct {
	Name     string
	F        func() error
	Duration time.Duration
	MaxTimes int
}

func NewTry(name string, f func() error, duration time.Duration, maxTimes int) *Try {
	return &Try{
		Name:     name,
		F:        f,
		Duration: duration,
		MaxTimes: maxTimes,
	}
}

type Report struct {
	Name        string        // 重试任务名称
	Result      bool          // 函数执行的结果
	Times       int           // 重试的次数
	SumDuration time.Duration // 总执行时间
	Errs        []error       // 函数执行的错误记录
}

func (r *Report) Error() string {
	return fmt.Sprintf("[retry]名称：%s，结果：%v，尝试次数：%v，总时间：%v，错误：%v", r.Name, r.Result, r.Times, r.SumDuration, r.Errs)
}

// Run 尝试重试，返回 chan 可以用于接收尝试报告
func (try *Try) Run() <-chan Report {
	result := make(chan Report, 1)
	go func() {
		defer close(result)
		start := time.Now()
		var errs []error
		for i := 0; i < try.MaxTimes; i++ {
			time.Sleep(try.Duration)
			err := try.F()
			if err == nil {
				result <- Report{
					Name:        try.Name,
					Result:      true,
					Times:       i + 1,
					SumDuration: time.Since(start),
					Errs:        errs,
				}
				return
			}
			errs = append(errs, err)
		}
		result <- Report{
			Name:        try.Name,
			Result:      false,
			Times:       try.MaxTimes,
			SumDuration: time.Since(start),
			Errs:        errs,
		}
	}()
	return result
}
