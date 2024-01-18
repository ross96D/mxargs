package execute

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"sync"
)

type print struct {
	mut *sync.Mutex
}

var p print = print{
	mut: &sync.Mutex{},
}

func (p *print) print(buff *bytes.Buffer) {
	p.mut.Lock()
	defer p.mut.Unlock()
	// buff := bytes.NewBuffer(b)
	io.Copy(os.Stdout, buff)
}

func Execute(cmd exec.Cmd, args []string) {
	for i := 0; i < len(args); i++ {
		cmd.Args = append(cmd.Args, args[i])
		go func(cmd exec.Cmd) {
			var b []byte
			buff := bytes.NewBuffer(b)
			r, _ := cmd.StdoutPipe()
			io.Copy(buff, r)
			r.Close()
			p.print(buff)
			cmd.Run()
		}(cmd)
	}
}
