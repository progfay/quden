package endpoint

import "fmt"

type Endpoint struct {
	Method string
	Path   string
	RegExp string
}

func New(method, path, regexp string) Endpoint {
	return Endpoint{
		Method: method,
		Path:   path,
		RegExp: regexp,
	}
}

func (e *Endpoint) String() string {
	return fmt.Sprintf("%s %s", e.Method, e.Path)
}
