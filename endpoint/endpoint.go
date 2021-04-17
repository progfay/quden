package endpoint

import "fmt"

type Endpoint struct {
	Method  string
	Pattern string
}

func New(method, pattern string) *Endpoint {
	return &Endpoint{
		Method:  method,
		Pattern: pattern,
	}
}

func (e *Endpoint) String() string {
	return fmt.Sprintf("%s %s", e.Method, e.Pattern)
}
