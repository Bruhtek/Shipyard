package utils

type EnvDescription struct {
	Name    string
	EnvType string

	// only applicable to remote environments
	EnvKey    string `json:"-"`
	Heartbeat bool
}
