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

	Output string
	ctx    context.Context

	Mutex sync.RWMutex `json:"-"`
}

type ActionStatus int

const (
	Pending ActionStatus = 0
	Running ActionStatus = 1
	Success ActionStatus = 2
	Failed  ActionStatus = 3
)
