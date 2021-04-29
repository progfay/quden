package gorilla

import (
	"testing"
)

func Test_parsePath(t *testing.T) {
	for _, testcase := range []struct {
		title string
		in    string
		want  struct {
			name    string
			pattern string
		}
	}{
		{
			title: "static path",
			in:    `/users`,
			want: struct {
				name    string
				pattern string
			}{
				name:    `/users`,
				pattern: `/users`,
			},
		},
		{
			title: "variable path (single)",
			in:    `/users/{id}`,
			want: struct {
				name    string
				pattern string
			}{
				name:    `/users/{id}`,
				pattern: `/users/[^/]+`,
			},
		},
		{
			title: "variable path with custom RegExp (single)",
			in:    `/users/{id:\d+}`,
			want: struct {
				name    string
				pattern string
			}{
				name:    `/users/{id}`,
				pattern: `/users/\d+`,
			},
		},
		{
			title: "variable path (multiple)",
			in:    `/users/{id}/files/{filename}.{ext}`,
			want: struct {
				name    string
				pattern string
			}{
				name:    `/users/{id}/files/{filename}.{ext}`,
				pattern: `/users/[^/]+/files/[^/]+\.[^/]+`,
			},
		},
		{
			title: "variable path with custom RegExp (multiple)",
			in:    `/users/{id:\d+}/files/{filename}.{ext:png|jpe?g|gip|pdf}`,
			want: struct {
				name    string
				pattern string
			}{
				name:    `/users/{id}/files/{filename}.{ext}`,
				pattern: `/users/\d+/files/[^/]+\.png|jpe?g|gip|pdf`,
			},
		},
	} {
		t.Run(testcase.title, func(t *testing.T) {
			name, pattern, err := parsePath(testcase.in)
			if err != nil {
				t.Error(err)
			}

			if name != testcase.want.name {
				t.Errorf("want name %q, got %q", testcase.want.name, name)
			}

			if pattern != testcase.want.pattern {
				t.Errorf("want pattern %q, got %q", testcase.want.pattern, pattern)
			}
		})
	}
}
