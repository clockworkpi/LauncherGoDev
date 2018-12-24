package UI

import (
  "sync"

)

type FolderStack struct {
    lock *sync.Mutex
    head *element
    Size int
    RootPath string
}

func (stk *FolderStack) Push(data interface{}) {
    stk.lock.Lock()

    element := new(element)
    element.data = data
    temp := stk.head
    element.next = temp
    stk.head = element
    stk.Size++

    stk.lock.Unlock()
}

func (stk *FolderStack) Pop() interface{} {
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

func (stk *FolderStack) SetRootPath(path string) {
  stk.RootPath = path
}

func (stk *FolderStack) Length() int {
	return stk.Size
}

func (stk *FolderStack) Last() string {
  idx := stk.Length() -1
  if idx < 0 {
    return stk.RootPath
  }else {
    return stk.head.data.(string)
  }
}

func NewFolderStack() *FolderStack {
    stk := new(FolderStack)
    stk.lock = &sync.Mutex{}
    return stk
}
