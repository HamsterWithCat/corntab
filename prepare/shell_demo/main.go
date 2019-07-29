package main

import (
	"fmt"
	"os/exec"
)

func main(){
	var(
		cmd *exec.Cmd
	)
	cmd = exec.Command("/bin/bash","-c","echo 1;echo 2;")
	err := cmd.Run()
	if err!=nil{
		fmt.Println("execute command err,err ",err.Error())
		return
	}
}
