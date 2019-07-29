package main

import (
	"context"
	"fmt"
	"time"
)

func main(){
	//时间格式化
	//fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	//context 控制
	ctx,cancelFunc := context.WithCancel(context.Background())

	go func(ctx context.Context){
		for{
			fmt.Println("hahaha")
			time.Sleep(time.Millisecond*500)
			select {
			case <-ctx.Done():
				fmt.Println("aaaaa")
				return
			default:

			}
		}
	}(ctx)

	time.Sleep(time.Second*2)
	cancelFunc()
	fmt.Println("----------cancel------------")
	time.Sleep(time.Second*5)
}
