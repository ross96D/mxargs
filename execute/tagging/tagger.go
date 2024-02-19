package execute

import (
	"bytes"
)

type tagged struct {
	Tag  int
	buff *bytes.Buffer
}

func New(tag int) tagged {
	buff := make([]byte, 0, 1000)
	return tagged{Tag: tag, buff: bytes.NewBuffer(buff)}
}

func (t *tagged) Read(p []byte) (n int, err error) {
	return t.buff.Read(p)
}

func (t *tagged) Write(p []byte) (n int, err error) {
	return t.buff.Write(p)
}
