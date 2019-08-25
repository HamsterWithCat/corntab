package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

//两个问题，协程之间如何通信
//子协程退出意味着任务结束吗 结束 子进程结束后，意味着任务放弃执行，暂时是这样

//对于context包，父进程控制子进程的结束后，子进程会返回进程退出的信息，而不是直接退出造成卡死
type Result struct {
	Result []byte
	Err    error
}

var resultChan chan *Result

func init() {
	//初始化channel
	resultChan = make(chan *Result, 1)
}
func main() {
	ctxRoot := context.Background()
	//context控制上下文
	ctxCancel, cancelFunc := context.WithCancel(ctxRoot)
	//在协程内执行command
	go func(ctx context.Context, resultChan chan<- *Result) {
		cmd := exec.CommandContext(ctxCancel, "/bin/bash", "-c", "ls /")
		fmt.Println("------cmd已创建------")
		time.Sleep(time.Second * 2)
		//定义结构体接收命令执行结果
		var result Result
		result.Result, result.Err = cmd.CombinedOutput()
		fmt.Println("----------命令执行完成----------")

		//执行结果写入channel
		resultChan <- &result
	}(ctxCancel, resultChan)

	//主线程sleep1秒后，取消任务的执行
	time.Sleep(time.Second * 3)
	cancelFunc()

	//主线程轮训result
	fmt.Println("------轮训开始---------")
	select {
	case result := <-resultChan:
		fmt.Printf("%T", result)
		fmt.Println(string(result.Result), "  ", result.Err)
	}
}
