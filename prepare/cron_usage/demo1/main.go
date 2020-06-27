package main

import (
	"github.com/gorhill/cronexpr"
	"fmt"
	"time"
)

func main(){
	var (
		expr *cronexpr.Expression
		err error
		now time.Time
		nextTime time.Time
	)
	//linux crontab
	//哪一分钟（0-59），哪一小时（0-23），哪一天（1-31），哪月（1-12），星期几（0-6）
	/*//每分钟执行1次
	if expr,err = cronexpr.Parse("* * * * *");err != nil {
		fmt.Println(err)
		return
	}*/
	//每隔5分钟执行1次
	//if expr,err = cronexpr.Parse("*/5 * * * *");err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//0, 5, 10, 15, 20...执行时间

	//github.com/gorhill/cronexpr这个包可以执行到秒
	//每5秒执行一次
	if expr,err = cronexpr.Parse("*/5 * * * * * *");err != nil {
		fmt.Println(err)
		return
	}

	//当前时间
	now = time.Now()
	//下次调度时间
	nextTime = expr.Next(now)
	fmt.Println(now,nextTime)

	//nextTime.Sub(now),下一次调度时间减去现在时间
	//等待这个定时器超时,回调函数会被回调
	time.AfterFunc(nextTime.Sub(now), func() {
		fmt.Println("被调度了：",nextTime)
	})

	time.Sleep(5 * time.Second)
}
