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

func ExecuteWithOrder(cmd []string, args []string) {
	w := sync.WaitGroup{}

	ctx := context.Background()
	g, _ := errgroup.WithContext(ctx)
	g.SetLimit(runtime.NumCPU())

	orderedStdout := ordering.New()
	orderedStderr := ordering.New()

	w.Add(1)
	go func() {
		io.Copy(os.Stdout, orderedStdout)
		w.Done()
	}()
	go func() {
		io.Copy(os.Stderr, orderedStderr)
		w.Done()
	}()

	N := len(args)
	for i := 0; i < N; i++ {
		w.Add(1)
		cmd2 := exec.Command(cmd[0], cmd[1:]...)
		cmd2.Args = append(cmd2.Args, args[i])

		g.Go(func(cmd exec.Cmd, index int) func() error {
			return func() error {
				// chain stdout
				b := make([]byte, 0, 1024)
				bufferOut := bytes.NewBuffer(b)
				mout := &BufferTreadSafe{buffer: bufferOut}
				orderedStdout.Add(index, mout, index == N-1)
				cmd.Stdout = mout

				// chain stderr
				b = make([]byte, 0, 1024)
				bufferErr := bytes.NewBuffer(b)
				merr := &BufferTreadSafe{buffer: bufferErr}
				orderedStderr.Add(index, merr, index == N-1)
				cmd.Stderr = bufferErr

				//! TODO handle error
				cmd.Run()
				mout.ended = true

				w.Done()
				return nil
			}
		}(*cmd2, i))
	}
	w.Wait()
}

type BufferTreadSafe struct {
	buffer *bytes.Buffer
	ended  bool

	m sync.Mutex
}

func (m *BufferTreadSafe) Write(p []byte) (int, error) {
	m.m.Lock()
	defer m.m.Unlock()

	return m.buffer.Write(p)
}

func (m *BufferTreadSafe) Read(p []byte) (int, error) {
	m.m.Lock()
	defer m.m.Unlock()

	n, err := m.buffer.Read(p)
	if n == 0 && err == io.EOF && !m.ended {
		err = nil
	}
	return n, err
}
