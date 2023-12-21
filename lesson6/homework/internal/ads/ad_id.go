package ads

import "sync"

func NewAdId() int64 {
	return ai.GenerateID()
}

var ai autoInc

type autoInc struct {
	sync.Mutex // ensures autoInc is goroutine-safe
	id         int64
}

func (a *autoInc) GenerateID() (id int64) {
	a.Lock()
	defer a.Unlock()

	id = a.id
	a.id++
	return
}
