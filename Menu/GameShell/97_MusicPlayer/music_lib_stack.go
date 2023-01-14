package MusicPlayer

import (
	"sync"
)

type element struct {
	data interface{}
	next *element
}

type MusicLibStack struct {
	lock *sync.Mutex
	head *element
	Size int
}

func (stk *MusicLibStack) Push(data interface{}) {
	stk.lock.Lock()

	element := new(element)
	element.data = data
	temp := stk.head
	element.next = temp
	stk.head = element
	stk.Size++

	stk.lock.Unlock()
}

func (stk *MusicLibStack) Pop() interface{} {
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

func (stk *MusicLibStack) Length() int {
	return stk.Size
}

func (stk *MusicLibStack) Last() string {
	idx := stk.Length() -1
	if idx < 0 {
		return "/"
	} else {
		return stk.head.data.(string)
	}
}

func NewMusicLibStack() *MusicLibStack {
	stk := new(MusicLibStack)
	stk.lock = &sync.Mutex{}
	return stk
}


