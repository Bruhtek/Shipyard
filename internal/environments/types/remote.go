package types

type RequestData struct {
	Path   string
	Method string
	Body   string
}

type RequestResponse struct {
	Body string
	Code int
}
