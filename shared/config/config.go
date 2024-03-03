package config

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/ross96D/mxargs/shared/flags"
	sysconf "github.com/tklauser/go-sysconf"
)

const (
	maxSize = 128 * 1024
)

type Configuration struct {
	Input io.Reader
	//! xargs made as default spaces and new lines... now this only handles one case.
	Delimeter byte
	sc        *bufio.Scanner
	text      string
	peeked    bool

	MaxArgs int

	MaxChars int64

	EOF string
}

func New(flags *flags.CmdFlags, input io.Reader) (conf *Configuration, err error) {
	conf = new(Configuration)

	if flags.FilePath != "" {
		var f *os.File
		f, err = os.Open(flags.FilePath)
		if err != nil {
			return
		}
		conf.Input = f
	} else {
		if input == nil {
			err = errors.New("new_config: no input reader supplied")
			return
		}
		conf.Input = input
	}

	conf.Delimeter = flags.Delimiter[0]
	if conf.Delimeter == 0 {
		conf.Delimeter = '\n'
	}
	if flags.NullSeparator {
		conf.Delimeter = 0
	}

	conf.MaxArgs = flags.MaxArgs

	if conf.MaxChars, err = maxChars(); err != nil {
		return
	}
	if flags.MaxChars != 0 {
		if flags.MaxChars > uint64(conf.MaxArgs) {
			err = fmt.Errorf("max number of character %d surpased", conf.MaxChars)
			return
		}
		conf.MaxChars = int64(flags.MaxChars)
	}

	// initialize the scanner
	conf.scanner()
	return
}

func (m *Configuration) splitFunc() bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		if i := bytes.IndexByte(data, m.Delimeter); i >= 0 {
			return i + 1, data[0:i], nil
		}

		// TODO check if this stop
		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return len(data), data, nil
		}
		return
	}
}

func (m *Configuration) scanner() {
	sc := bufio.NewScanner(m.Input)
	sc.Split(m.splitFunc())
	m.sc = sc
}

func (m *Configuration) scan() bool {
	//!!!! TODO this need refactor.......
	if m.peeked {
		m.peeked = false
		return true
	}
	next := m.sc.Scan()
	text := m.sc.Text()
	if m.EOF != "" && m.text == m.EOF {
		return false
	}
	m.text = text
	return next
}

func (m *Configuration) HasMoreData() bool {
	//!!!! TODO this need refactor.......
	resp := m.scan()
	m.peeked = true
	return resp
}

func (m *Configuration) Text() string {
	text := m.text
	m.text = ""
	return text
}

func (m *Configuration) Args() ([]string, error) {
	//! cannot use MaxArgs because it could be a number too big
	resp := make([]string, 0)
	for range m.MaxArgs {
		if m.scan() {
			resp = append(resp, m.Text())
		} else {
			text := m.Text()
			if text != "" {
				resp = append(resp, text)
			}
			return resp, io.EOF
		}
	}
	return resp, nil
}

var configuration Configuration

func Config() *Configuration {
	return &configuration
}

func maxChars() (int64, error) {
	max, err := sysconf.Sysconf(sysconf.SC_ARG_MAX)

	if err != nil {
		return 0, err
	}

	var envsize int64 = 0
	envs := os.Environ()
	for _, env := range envs {
		envsize += int64(len(env)) + 1
	}

	result := (max - envsize - 2048)
	if result > maxSize {
		result = maxSize
	}

	return result, nil
}
