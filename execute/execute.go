package execute

import (
	"bytes"
	"context"
	"fmt"
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

func (p *Print) print(m *MWriter) error {
	p.Mut.Lock()
	defer p.Mut.Unlock()
	buff := bytes.NewBuffer(m.buff)

	_, err := io.Copy(os.Stdout, buff)
	if err != nil {
		return err
	}
	return nil
}

func Execute(print *Print, cmd []string, args []string) {
	w := sync.WaitGroup{}

	ctx := context.Background()
	g, _ := errgroup.WithContext(ctx)
	g.SetLimit(runtime.NumCPU())

	for i := 0; i < len(args); i++ {
		w.Add(1)
		cmd2 := exec.Command(cmd[0], cmd[1:]...)
		cmd2.Args = append(cmd2.Args, args[i])

		g.Go(func(cmd exec.Cmd) func() error {
			return func() error {
				m := &MWriter{}
				cmd.Stdout = m

				cmd.Run()
				if err := print.print(m); err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
				}
				w.Done()
				return nil
			}
		}(*cmd2))
	}
	w.Wait()
}

type MWriter struct {
	buff []byte
}

func (m *MWriter) Write(p []byte) (int, error) {
	if m.buff == nil {
		m.buff = make([]byte, 0, len(p))
	}
	m.buff = append(m.buff, p...)
	return len(p), nil
}
