package UI

import (
	"sync"
)

type element struct {
        data interface{}
        next *element
}

type PageStack struct {
	lock     *sync.Mutex
	head     *element
	Size     int
}

func (stk *PageStack) Push(data interface{}) {
	stk.lock.Lock()

	element := new(element)
	element.data = data
	temp := stk.head
	element.next = temp
	stk.head = element
	stk.Size++

	stk.lock.Unlock()
}

func (stk *PageStack) Pop() interface{} {
	if stk.head == nil {
		return nil
	}
	stk.lock.Lock()
	r := stk.head.data
	stk.head = stk.head.next
	stk.Size--

	stk.lock.Unlock()

	return r
}

func (stk *PageStack) Length() int {
	return stk.Size
}

func (stk *PageStack) Last() interface{} {
	idx := stk.Length() - 1
	if idx < 0 {
		return nil
	} else {
		return stk.head.data
	}
}

func NewPageStack() *PageStack {
	stk := new(PageStack)
	stk.lock = &sync.Mutex{}
	return stk
}
