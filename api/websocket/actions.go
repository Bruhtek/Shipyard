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

	Command    []string
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

func (a *Action) Retry() (res bool) {
	defer func() {
		if r := recover(); r != nil {
			println("Panic while retrying action:")
			println(r)
			a.Mutex.Unlock()
			res = false
		}
	}()

	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	if a.Status != Failed {
		return false
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	runner := Runner{
		Command:  a.Command,
		ActionId: a.ActionId,
		Action:   a,
		Ctx:      ctx,
	}
	a.ctx = ctx
	a.cancelFunc = cancelFunc

	a.Status = Running
	a.StartedAt = time.Now()
	a.FinishedAt = time.Time{}
	a.Output += "\x1b[2J\x1b[H" // this clear the screen
	ConnectionManager.BroadcastActionOutput(a.ActionId, "\x1b[2J\x1b[H")

	go runner.Run()

	return true
}

type ActionStatus int

const (
	Pending ActionStatus = 0
	Running ActionStatus = 1
	Success ActionStatus = 2
	Failed  ActionStatus = 3
)
