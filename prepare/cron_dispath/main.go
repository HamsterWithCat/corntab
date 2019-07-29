package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

type CronJob struct {
	expr *cronexpr.Expression
	nextTime time.Time
}
var scheduleTable map[string]*CronJob
func init(){
	//初始化调度表
	scheduleTable = make(map[string]*CronJob)
}
func main(){
	//一个调度协程，定时检查所有的Cron任务，执行过期任务

	//定义两个定时任务
	expr := cronexpr.MustParse("*/5 * * * * * *")
	job1 := &CronJob{
		expr:expr,
		nextTime:expr.Next(time.Now()),
	}
	//将任务添加到任务调度表中
	scheduleTable["job1"]=job1

	//定义两个定时任务
	expr = cronexpr.MustParse("*/6 * * * * * *")
	job2 := &CronJob{
		expr:expr,
		nextTime:expr.Next(time.Now()),
	}
	//将任务添加到任务调度表中
	scheduleTable["job2"]=job2

	//启动一个协程对调度表进行轮训
	go func(){
		//循环检查
		for{
			//获取当前时间
			now := time.Now()
			//轮训所有job
			for jobName,job :=range scheduleTable{
				//判断job是否到期
				if !job.nextTime.After(now){
					//job到期 启动协程执行任务
					fmt.Println(jobName+"已经到期")
					go func(cronJob *CronJob){
						fmt.Println("正在执行job")
						time.Sleep(time.Millisecond*100)//模拟执行过程
						//执行完成 修改job下次到期的时间
						cronJob.nextTime = cronJob.expr.Next(time.Now())
					}(job)
				}
			}
			//轮训协程休眠  减少空轮训次数
			select {
				case <- time.NewTimer(time.Millisecond*1000).C:
			}
		}
	}()

	//主线程等待协程执行
	time.Sleep(time.Second*10)
}
