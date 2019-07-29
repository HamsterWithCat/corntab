package main

import(
	"log"
	"github.com/gorhill/cronexpr"
	"time"
)
func main(){
	expr,err := cronexpr.Parse("*/2 * * * * * *")
	if err != nil {
		log.Printf("parse cron expression failed,err %s",err.Error())
		return
	}
	nextTime := expr.Next(time.Now())
	log.Printf("%v\n",nextTime)

	//fmt.Println(time.Now()," ",time.Now().Unix()," ",time.Now().Nanosecond())

	time.AfterFunc(nextTime.Sub(time.Now()), func() {
		log.Printf("hello\n")
		time.Sleep(time.Second*5)
	})

	//主线程空转等待子线程执行
	select {

	}//select 空转时如果所有的子协程已经退出，会报死锁错误
	//time.Sleep(time.Second*10)
}
