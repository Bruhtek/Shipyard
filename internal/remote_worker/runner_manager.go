package remote_worker

import (
	"Shipyard/internal/terminals"
	"Shipyard/internal/utils"
	"context"
	"sync"
)

type RunnerMngr struct {
	runners map[string]*terminals.Runner

	runnersMutex sync.RWMutex
}

var RunnerManager RunnerMngr = RunnerMngr{
	runners:      make(map[string]*terminals.Runner),
	runnersMutex: sync.RWMutex{},
}

func (r *RunnerMngr) NewRunner(id string, cmd []string) *terminals.Runner {
	r.runnersMutex.Lock()
	defer r.runnersMutex.Unlock()

	ctx, cancel := context.WithCancel(context.Background())

	runner := terminals.Runner{
		Command:    cmd,
		Ctx:        ctx,
		CancelFunc: cancel,

		OutputFn: func(out string) {
			message := map[string]interface{}{
				"Output": out,
				"Type":   "OutputFn",
			}
			CManager.SendResponse(id, message)
		},
		OutputMetaFn: func(status utils.ActionStatus) {
			message := map[string]interface{}{
				"ActionStatus": status,
				"Type":         "OutputMetaFn",
			}
			CManager.SendResponse(id, message)
		},
		DeleteFn: func() {
			message := map[string]interface{}{
				"Type": "DeleteFn",
			}
			CManager.SendResponse(id, message)
		},
	}

	go runner.Run()

	r.runners[id] = &runner

	return &runner
}

func (r *RunnerMngr) RetryRunner(id string) {
	r.runnersMutex.Lock()
	defer r.runnersMutex.Unlock()

	runner, ok := r.runners[id]
	if !ok {
		return
	}

	runner.CancelFunc()

	ctx, cancelFunc := context.WithCancel(context.Background())
	runner.Ctx = ctx
	runner.CancelFunc = cancelFunc

	go runner.Run()
}

func (r *RunnerMngr) CancelRunner(id string) {
	r.runnersMutex.Lock()
	defer r.runnersMutex.Unlock()

	if runner, ok := r.runners[id]; ok {
		runner.OutputMetaFn(utils.Failed)
		runner.CancelFunc()
	}
}

func (r *RunnerMngr) DeleteRunner(id string) {
	r.runnersMutex.Lock()
	defer r.runnersMutex.Unlock()

	if runner, ok := r.runners[id]; ok {
		runner.CancelFunc()
		runner.DeleteFn()
	}

	delete(r.runners, id)
}
