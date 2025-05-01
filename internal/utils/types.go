package utils

type EnvDescription struct {
	Name    string
	EnvType string
}

type ActionStatus int

const (
	Pending ActionStatus = 0
	Running ActionStatus = 1
	Success ActionStatus = 2
	Failed  ActionStatus = 3
)
