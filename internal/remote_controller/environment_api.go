package remote_controller

func (r *RemoteEnvironment) GetResponse(path string) ([]byte, error) {
	r.Need()

	return []byte("{}"), nil
}
func (r *RemoteEnvironment) PostResponse(path string, body string) ([]byte, error) {
	r.Need()

	return []byte("{}"), nil
}
