package gorilla

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

// Ref. https://github.com/gorilla/mux/blob/98cb6bf42e086f6af920b965c38cacc07402d51b/regexp.go#L282-L310
// Ref. https://github.com/gorilla/mux/blob/98cb6bf42e086f6af920b965c38cacc07402d51b/regexp.go#L31-L150
func parsePath(path string) (string, string, error) {
	var name, pattern bytes.Buffer
	var level, idx int

	for i := 0; i < len(path); i++ {
		switch path[i] {
		case '{':
			if level++; level == 1 {
				name.WriteString(path[idx:i])
				pattern.WriteString(regexp.QuoteMeta(path[idx:i]))
				idx = i
			}
		case '}':
			if level--; level == 0 {
				parts := strings.SplitN(path[idx+1:i], ":", 2)
				pat := "[^/]+"
				if len(parts) == 2 {
					pat = parts[1]
				}
				name.WriteString(fmt.Sprintf("{%s}", parts[0]))
				pattern.WriteString(pat)
				idx = i + 1
			} else if level < 0 {
				return "", "", fmt.Errorf("unbalanced braces in %q", path)
			}
		}
	}

	if level != 0 {
		return "", "", fmt.Errorf("unbalanced braces in %q", path)
	}

	name.WriteString(path[idx:])
	pattern.WriteString(regexp.QuoteMeta(path[idx:]))
	return name.String(), pattern.String(), nil
}
