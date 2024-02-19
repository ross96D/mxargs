package execute

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sync"

	"github.com/ross96D/mxargs/execute/ordering"
	"golang.org/x/sync/errgroup"
)

func ExecuteWithOrder(print *Print, cmd []string, args []string) {
	w := sync.WaitGroup{}

	ctx := context.Background()
	g, _ := errgroup.WithContext(ctx)
	g.SetLimit(runtime.NumCPU())

	orderedStdout := ordering.New()
	orderedStderr := ordering.New()
	go io.Copy(os.Stdout, orderedStdout)
	go io.Copy(os.Stderr, orderedStderr)

	N := len(args)
	for i := 0; i < N; i++ {
		w.Add(1)
		cmd2 := exec.Command(cmd[0], cmd[1:]...)
		cmd2.Args = append(cmd2.Args, args[i])

		g.Go(func(cmd exec.Cmd, index int) func() error {
			return func() error {
				// chain stdout
				buffStdout := make([]byte, 0, 500)
				mstdout := bytes.NewBuffer(buffStdout)
				orderedStdout.Add(index, mstdout, index == N-1)
				cmd.Stdout = mstdout

				// chain stderr
				buffStderr := make([]byte, 0, 500)
				mstderr := bytes.NewBuffer(buffStderr)
				orderedStderr.Add(index, mstderr, index == N-1)
				cmd.Stderr = mstderr

				//! TODO handle error
				cmd.Run()

				w.Done()
				return nil
			}
		}(*cmd2, i))
	}
	w.Wait()
}
