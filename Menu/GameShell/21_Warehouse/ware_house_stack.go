package Warehouse

import (
	"sync"
)

type element struct {
	data interface{}
	next *element
}

type WareHouseStack struct {
	lock *sync.Mutex
	head *element
	Size int
}

func (stk *WareHouseStack) Push(data interface{}) {
	stk.lock.Lock()

	element := new(element)
	element.data = data
	temp := stk.head
	element.next = temp
	stk.head = element
	stk.Size++

	stk.lock.Unlock()
}

func (stk *WareHouseStack) Pop() interface{} {
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

func (stk *WareHouseStack) Length() int {
	return stk.Size
}

func (stk *WareHouseStack) Last() interface{} {
	idx := stk.Length() -1
	if idx < 0 {
		return nil
	} else {
		return stk.head.data
	}
}

func NewWareHouseStack() *WareHouseStack {
	stk := new(WareHouseStack)
	stk.lock = &sync.Mutex{}
	return stk
}


