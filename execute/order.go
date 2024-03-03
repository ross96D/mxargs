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
	"github.com/ross96D/mxargs/shared/config"
	"golang.org/x/sync/errgroup"
)

// ! test if passing a command (not a reference, a value, to avoid concurrency problems) insted of the slice of string works
func ExecuteWithOrder(cmd []string, conf *config.Configuration) {
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

	var isLast bool
	for i := 0; true; i++ {
		w.Add(1)
		command := exec.Command(cmd[0], cmd[1:]...)
		args, err := conf.Args()
		isLast = !conf.HasMoreData()
		if err != nil {
			isLast = true
			if len(args) == 0 {
				break
			}
		}
		command.Args = append(command.Args, args...)

		g.Go(func(cmd exec.Cmd, index int, isLast bool) func() error {
			return func() error {
				// chain stdout
				b := make([]byte, 0, 1024)
				bufferOut := bytes.NewBuffer(b)
				mout := &BufferTreadSafe{buffer: bufferOut}
				orderedStdout.Add(index, mout, isLast)
				cmd.Stdout = mout

				// chain stderr
				b = make([]byte, 0, 1024)
				bufferErr := bytes.NewBuffer(b)
				merr := &BufferTreadSafe{buffer: bufferErr}
				orderedStderr.Add(index, merr, isLast)
				cmd.Stderr = bufferErr

				//! TODO handle error
				cmd.Run()
				mout.ended = true

				w.Done()
				return nil
			}
		}(*command, i, isLast))

		if isLast {
			break
		}
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
