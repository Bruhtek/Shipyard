package websocket

import (
	"Shipyard/internal/env_manager"
	"Shipyard/internal/logger"
	"Shipyard/internal/terminals"
	"Shipyard/internal/utils"
	"context"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type Action struct {
	// Action description:
	Environment string
	IsRemote    bool
	Object      string
	Action      string
	ObjectId    string // optional

	// Metadata
	ActionId      string
	InitializedBy string // websocket connection id TODO: change that to user id
	StartedAt     time.Time
	FinishedAt    time.Time
	Status        utils.ActionStatus

	Command    []string
	Output     string
	Ctx        context.Context    `json:"-"`
	CancelFunc context.CancelFunc `json:"-"`

	Mutex sync.RWMutex `json:"-"`

	Broadcaster *Broadcaster `json:"-"`
}

type Broadcaster struct {
	BroadcastFn     func(string, interface{})
	BroadcastMetaFn func(action *Action)
	BroadcastMiscFn func(actionId string, key string, msg interface{})
}

func NewBroadcastAction(cmd []string, broadcaster *Broadcaster, env string, obj string, act string, objId string) *Action {
	action := newAction(cmd, env, obj, act, objId)
	action.Broadcaster = broadcaster
	return action
}

func newAction(cmd []string, env string, obj string, act string, objId string) *Action {
	actionId := utils.RandString(32)
	ctx, cancelFunc := context.WithCancel(context.Background())

	return &Action{
		Environment:   env,
		Object:        obj,
		Action:        act,
		ObjectId:      objId,
		ActionId:      actionId,
		InitializedBy: "",
		StartedAt:     time.Now(),
		FinishedAt:    time.Time{},
		Status:        utils.Pending,
		Output:        "",
		Command:       cmd,
		Ctx:           ctx,
		CancelFunc:    cancelFunc,
	}
}

func (a *Action) Cancel() (res bool) {
	defer func() {
		if r := recover(); r != nil {
			logger.HandleSimpleRecoverPanic(r, "Panic while cancelling action")
			a.Mutex.Unlock()
			res = false
		}
	}()

	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	a.CancelFunc()

	if a.Status == utils.Running || a.Status == utils.Pending {
		a.Status = utils.Failed
		a.FinishedAt = time.Now()
	}

	return true
}

func (a *Action) Retry() (res bool) {
	defer func() {
		if r := recover(); r != nil {
			logger.HandleSimpleRecoverPanic(r, "Panic while retrying action")
			a.Mutex.Unlock()
			res = false
		}
	}()

	a.Mutex.Lock()
	if a.Status != utils.Failed {
		a.Mutex.Unlock()
		log.Warn().
			Strs("cmd", a.Command).
			Str("action", a.Action).
			Msg("Trying to retry a non-failed action")
		return false
	}
	a.Ctx, a.CancelFunc = context.WithCancel(context.Background())
	a.Mutex.Unlock()

	a.Mutex.RLock()
	if a.IsRemote {
		remoteEnv := env_manager.EnvManager.GetEnv(a.Environment)
		if remoteEnv == nil || remoteEnv.GetEnvType() != "remote" {
			return false
		}
		env := remoteEnv.(env_manager.RemoteEnvironment)

		runner := terminals.RemoteRunner{
			Command:      a.Command,
			Ctx:          a.Ctx,
			CancelFunc:   a.CancelFunc,
			ID:           a.ActionId,
			Env:          env,
			OutputFn:     a.HandleOutput,
			OutputMetaFn: a.HandleMetadata,
			DeleteFn:     a.HandleDelete,
		}
		go runner.Retry()
	} else {
		runner := terminals.Runner{
			Ctx:          a.Ctx,
			Command:      a.Command,
			OutputFn:     a.HandleOutput,
			OutputMetaFn: a.HandleMetadata,
			DeleteFn:     a.HandleDelete,
		}
		go runner.Run()
	}
	a.Mutex.RUnlock()

	a.HandleOutput("\x1b[2J\x1b[H") // this clear the screen

	return true
}

func (a *Action) HandleOutput(output string) {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	a.Output += output
	if a.Broadcaster != nil {
		a.Broadcaster.BroadcastFn(a.ActionId, output)
	}
}

func (a *Action) HandleMetadata(status utils.ActionStatus) {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	a.Status = status
	if status == utils.Success || status == utils.Failed {
		a.FinishedAt = time.Now()
	}
	if status == utils.Running || status == utils.Pending {
		a.StartedAt = time.Now()
		a.FinishedAt = time.Time{}
	}

	if a.Broadcaster != nil {
		a.Broadcaster.BroadcastMetaFn(a)
	}
}

func (a *Action) HandleDelete() {
	if a.Broadcaster != nil {
		time.Sleep(10 * time.Second)
		a.Broadcaster.BroadcastMiscFn(a.ActionId, "Removed", true)
		ActionManager.DeleteFinishedAction(a, 500*time.Millisecond)
	}
}
