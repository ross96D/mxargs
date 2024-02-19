package ordering

import (
	"io"
	"sync"
)

type Ordering struct {
	queue    PriorityQueue[io.Reader]
	_current *Item[io.Reader]

	mut     sync.Mutex
	readMut sync.Mutex

	// this marks that all elements have been processing and will not add more elements
	end bool

	// used to wait for an element if the queue is empty there is still element to be added
	wait sync.Cond
}

func New() *Ordering {
	condMutex := sync.Mutex{}
	return &Ordering{
		queue: make(PriorityQueue[io.Reader], 0),
		mut:   sync.Mutex{},
		wait:  *sync.NewCond(&condMutex),
	}
}

func (o *Ordering) Add(tag int, reader io.Reader, last bool) {
	if o.end {
		return
	}

	o.mut.Lock()
	defer o.mut.Unlock()

	o.end = last

	item := NewItem(reader, tag)
	o.queue.Push(item)
	o.wait.Signal()
}

func (o *Ordering) setCurrent(current *Item[io.Reader]) {
	o.mut.Lock()
	defer o.mut.Unlock()

	o._current = current
}

func (o *Ordering) pop() (Item[io.Reader], error) {
	o.mut.Lock()
	defer o.mut.Unlock()

	if o.queue.Len() == 0 {
		if o.end {
			return Item[io.Reader]{}, io.EOF
		}
		o.mut.Unlock()

		// waits for the next element
		o.wait.L.Lock()
		o.wait.Wait()
		o.wait.L.Unlock()

		o.mut.Lock()
	}
	item := o.queue.Pop().(Item[io.Reader])
	// o.mut.Unlock()
	return item, nil
}

// i dont use a mutex on the read function because i asume this will not be read by several elements at the same time.
// so reading several times is an error
func (o *Ordering) Read(b []byte) (n int, err error) {
	if o._current == nil {
		current, err := o.pop()
		// if there is an error here is that all elements have been processed
		if err != nil {
			return 0, err
		}
		o.setCurrent(&current)
	}
	n, err = o._current.Value().Read(b)

	if err == io.EOF {
		o.setCurrent(nil)
		current, errPop := o.pop()
		// if no errors here means there is another element to process
		if errPop == nil {
			err = nil
			o.setCurrent(&current)
		}
	}
	return n, err
}
