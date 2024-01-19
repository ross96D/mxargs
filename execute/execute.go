package execute

import (
	"context"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sync"

	"golang.org/x/sync/errgroup"
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

	ctx := context.Background()
	g, _ := errgroup.WithContext(ctx)
	g.SetLimit(runtime.NumCPU())

	for i := 0; i < len(args); i++ {
		w.Add(1)
		cmd2 := *cmd
		cmd2.Args = append(cmd2.Args, args[i])

		g.Go(func(cmd exec.Cmd) func() error {
			return func() error {
				r, _ := cmd.StdoutPipe()
				go func() {
					cmd.Run()
					r.Close()
				}()
				print.print(r)
				w.Done()
				return nil
			}
		}(cmd2))

		// go func(cmd exec.Cmd) {

		// }(cmd2)
	}
	w.Wait()
}
