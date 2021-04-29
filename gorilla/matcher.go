package gorilla

import (
	"strings"
)

type artifact struct {
	path           string
	pathTerminated bool
	methodSet      map[string]struct{}
	isInvalid      bool
}

func newArtifact() *artifact {
	return &artifact{
		path:           "/",
		pathTerminated: false,
		methodSet:      nil,
		isInvalid:      false,
	}
}

type matcher interface {
	Process(art *artifact)
}

type pathMatcher struct {
	path string
}

func newPathMatcher(path string) pathMatcher {
	return pathMatcher{
		path: path,
	}
}

func (matcher pathMatcher) Process(art *artifact) {
	if art.isInvalid || art.pathTerminated {
		art.isInvalid = true
		return
	}

	art.path = strings.TrimRight(art.path, "/") + strings.TrimRight(matcher.path, "/")
	if art.path == "" {
		art.path = "/"
	}
	art.pathTerminated = true
}

type pathPrefixMatcher struct {
	pathPrefix string
}

func newPathPrefixMatcher(pathPrefix string) pathPrefixMatcher {
	return pathPrefixMatcher{
		pathPrefix: pathPrefix,
	}
}

func (matcher pathPrefixMatcher) Process(art *artifact) {
	if art.isInvalid || art.pathTerminated {
		art.isInvalid = true
		return
	}

	art.path = strings.TrimRight(art.path, "/") + strings.TrimRight(matcher.pathPrefix, "/")
	if art.path == "" {
		art.path = "/"
	}
}

type methodsMatcher struct {
	methodSet map[string]struct{}
}

func newMethodsMatcher(methodList []string) methodsMatcher {
	methodSet := map[string]struct{}{}
	for _, method := range methodList {
		methodSet[strings.ToUpper(method)] = struct{}{}
	}

	return methodsMatcher{
		methodSet: methodSet,
	}
}

func (matcher methodsMatcher) Process(art *artifact) {
	if art.isInvalid {
		return
	}

	if art.methodSet == nil {
		art.methodSet = matcher.methodSet
		art.isInvalid = len(art.methodSet) == 0
		return
	}

	for method := range art.methodSet {
		if _, ok := matcher.methodSet[method]; !ok {
			delete(art.methodSet, method)
		}
	}

	art.isInvalid = len(art.methodSet) == 0
}
