/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"io"
	"os"
	"os/exec"

	"github.com/ross96D/mxargs/execute"
)

func main() {
	// cmd.Execute()
	var mcmd *exec.Cmd
	args := os.Args
	args = args[1:]
	if len(args) > 1 {
		mcmd = exec.Command(args[0], args[1:]...)
	} else {
		mcmd = exec.Command(args[0])
	}
	_ = mcmd

	// var p = execute.Print{
	// 	Mut: &sync.Mutex{},
	// }
	argsss := stdin()
	execute.ExecuteWithOrder(args, argsss)
}

func stdin() []string {
	stdin := os.Stdin
	b, err := io.ReadAll(stdin)
	if err != nil {
		return []string{}
	}
	resp := make([]string, 0)

	buff := make([]byte, 0, len(b))
	for i := 0; i < len(b); i++ {
		if b[i] == 0 {
			resp = append(resp, string(buff))
			buff = buff[0:0]
			continue
		}
		buff = append(buff, b[i])
	}
	return resp
}
