package websocket

import (
	"sync"
	"time"
)

type AMStruct struct {
	actions map[string]*Action
	mutex   sync.RWMutex
}

func newActionManager() *AMStruct {
	return &AMStruct{
		actions: make(map[string]*Action),
		mutex:   sync.RWMutex{},
	}
}

var ActionManager = newActionManager()

func (am *AMStruct) createAction(action *Action) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	am.actions[action.ActionId] = action
}

func (am *AMStruct) GetAction(actionId string) (*Action, bool) {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	action, ok := am.actions[actionId]
	return action, ok
}

func (am *AMStruct) GetActions() map[string]*Action {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	return am.actions
}

func (am *AMStruct) GetEnvActions(envName string) map[string]*Action {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	envActions := make(map[string]*Action)
	for _, action := range am.actions {
		if action.Environment == envName {
			envActions[action.ActionId] = action
		}
	}
	return envActions
}

func (am *AMStruct) DeleteFinishedAction(action *Action, delay time.Duration) {
	if delay > 0 {
		time.Sleep(delay)
	}

	am.mutex.Lock()
	defer am.mutex.Unlock()

	action.Mutex.Lock()
	defer action.Mutex.Unlock()

	ConnectionManager.BroadcastActionMisc(action.ActionId, "Removed", true)

	am.actions[action.ActionId] = nil
	delete(am.actions, action.ActionId)
}
