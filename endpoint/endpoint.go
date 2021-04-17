package endpoint

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
