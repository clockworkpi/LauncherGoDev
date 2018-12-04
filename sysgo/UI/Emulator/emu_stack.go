package Emulator

import (
  "sync"
  "github.com/cuu/LauncherGoDev/sysgo/UI"

)

type element struct {
    data interface{}
    next *element
}

type EmuStack struct {
    lock *sync.Mutex
    head *element
    Size int
    EmulatorConfig *ActionConfig
}

func (stk *EmuStack) Push(data interface{}) {
    stk.lock.Lock()

    element := new(element)
    element.data = data
    temp := stk.head
    element.next = temp
    stk.head = element
    stk.Size++

    stk.lock.Unlock()
}

func (stk *EmuStack) Pop() interface{} {
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

func (stk *EmuStack) Length() int {
	return stk.Size
}

func (stk *EmuStack) Last() string {
  idx := stk.Length() -1
  if idx < 0 {
    return stk.EmulatorConfig.ROM
  }else {
    return stk.head.data.(string)
  }
}

func NewEmuStack() *EmuStack {
    stk := new(EmuStack)
    stk.lock = &sync.Mutex{}
    return stk
}
