package websocket

import (
	"context"
	"sync"
	"time"
)

type Action struct {
	// Action description:
	Environment string
	Object      string
	Action      string
	ObjectId    string // optional

	// Metadata
	ActionId      string
	InitializedBy string // websocket connection id TODO: change that to user id
	StartedAt     time.Time
	FinishedAt    time.Time
	Status        ActionStatus

	Output     string
	ctx        context.Context
	cancelFunc context.CancelFunc

	Mutex sync.RWMutex `json:"-"`
}

func (a *Action) Cancel() (res bool) {
	defer func() {
		if r := recover(); r != nil {
			println("Panic while cancelling action:")
			println(r)
			a.Mutex.Unlock()
			res = false
		}
	}()

	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	a.cancelFunc()

	if a.Status == Running || a.Status == Pending {
		a.Status = Failed
		a.FinishedAt = time.Now()
		return true
	}

	return true
}

type ActionStatus int

const (
	Pending ActionStatus = 0
	Running ActionStatus = 1
	Success ActionStatus = 2
	Failed  ActionStatus = 3
)
