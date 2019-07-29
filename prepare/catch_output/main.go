package main

import (
	"fmt"
	"os/exec"
)

func main(){
	//生成命令
	cmd := exec.Command("/bin/bash","-c","ls /")
	//执行命令，获取子进程的输出
	output,err := cmd.CombinedOutput();
	if err !=nil{
		fmt.Println("catch output failed,err",err.Error())
		return
	}
	fmt.Println(string(output))
}
