package parawri

import (
	"fmt"
	"io"
	"sync"

	ansi "github.com/k0kubun/go-ansi"
)

type Parallel struct {
	w   io.Writer
	cnt int
	mu  sync.Mutex
}

func NewParallelStdout() *Parallel {
	return &Parallel{w: ansi.NewAnsiStdout()}
}

func NewParallelStderr() *Parallel {
	return &Parallel{w: ansi.NewAnsiStderr()}
}

func (p *Parallel) NewAppendWriter() io.Writer {
	return &Writer{p: p, a: true}
}

func (p *Parallel) NewWriter() io.Writer {
	return &Writer{p: p, a: false}
}

type Writer struct {
	p   *Parallel
	idx int
	str string
	a   bool
}

func (w *Writer) Write(p []byte) (int, error) {
	w.p.mu.Lock()
	defer w.p.mu.Unlock()

	if w.a {
		w.str = w.str + string(p)
	} else {
		w.str = string(p)
	}

	if w.idx == 0 {
		w.p.cnt = w.p.cnt + 1
		w.idx = w.p.cnt
		fmt.Fprintln(w.p.w, w.str)
		return len(p), nil
	}

	diff := w.p.cnt - w.idx + 1
	ansi.CursorPreviousLine(diff)
	ansi.EraseInLine(1)
	fmt.Fprint(w.p.w, w.str)
	ansi.CursorNextLine(diff)
	return len(p), nil
}
