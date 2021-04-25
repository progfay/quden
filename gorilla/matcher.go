package gorilla

type matcherMode int8

const (
	pathMode matcherMode = iota
	pathPrefixMode
	methodsMode
)

type matcher struct {
	mode    matcherMode
	pattern string
	methods []string
}

func newPathMatcher(pattern string) matcher {
	return matcher{
		mode:    pathMode,
		pattern: pattern,
	}
}

func newPathPrefixMatcher(pattern string) matcher {
	return matcher{
		mode:    pathPrefixMode,
		pattern: pattern,
	}
}

func newMethodsMatcher(methods []string) matcher {
	return matcher{
		mode:    methodsMode,
		methods: methods,
	}
}
