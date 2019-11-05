package response

type SResponse struct {
	Byte       []byte
	Status     string
	StatusCode int
	Html       string
	Cookies    []string
	Location   string
}
