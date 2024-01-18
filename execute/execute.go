package execute

import (
	"io"
	"os"
	"os/exec"
	"sync"
)

type Print struct {
	Mut *sync.Mutex
}

func (p *Print) print(r io.Reader) {
	p.Mut.Lock()
	io.Copy(os.Stdout, r)
	p.Mut.Unlock()
}

func Execute(print *Print, cmd *exec.Cmd, args []string) {
	w := sync.WaitGroup{}
	for i := 0; i < len(args); i++ {
		w.Add(1)
		cmd2 := *cmd
		cmd2.Args = append(cmd2.Args, args[i])
		go func(cmd exec.Cmd) {
			r, _ := cmd.StdoutPipe()
			go func() {
				cmd.Run()
				r.Close()
			}()
			print.print(r)
			w.Done()
		}(cmd2)
	}
	w.Wait()
}
