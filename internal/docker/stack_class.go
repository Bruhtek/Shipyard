package docker

type Stack struct {
	Name        string
	Status      string
	ConfigFiles string

	Containers []*Container
	Networks   []*Network

	UpToDate UpToDateStatus
}
