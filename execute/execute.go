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

func Execute(cmd *exec.Cmd, args []string) {
	w := sync.WaitGroup{}
	for i := 0; i < len(args); i++ {
		w.Add(1)
		cmd2 := *cmd
		cmd2.Args = append(cmd2.Args, args[i])
		go func(cmd exec.Cmd) {
			// println("ZZ", cmd.String())
			var b []byte
			buff := bytes.NewBuffer(b)
			r, _ := cmd.StdoutPipe()
			go io.Copy(buff, r)
			cmd.Run()
			r.Close()
			p.print(buff)
			w.Done()
		}(cmd2)
	}
	w.Wait()
}
