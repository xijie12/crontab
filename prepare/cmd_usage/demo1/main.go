package main

import (
	"os/exec"
	"fmt"
)

func main(){
	var (
		cmd *exec.Cmd
		err error
	)

	//cmd = exec.Command("/bin/bash","-c","echo 1;echo 2;")

	cmd = exec.Command("/bin/bash","-c","echo 1")
	err = cmd.Run()
	fmt.Println(err)
}
