package goji_test

import (
	"fmt"
	"go/parser"
	"reflect"
	"testing"

	"github.com/progfay/quden/endpoint"
	"github.com/progfay/quden/goji"
)

func Test_MatchImportPath(t *testing.T) {
	framework := goji.New()

	for _, testcase := range []struct {
		in   string
		want bool
	}{
		{
			in:   "goji.io/pat",
			want: true,
		},
		{
			in:   "fmt",
			want: false,
		},
	} {
		t.Run(fmt.Sprintf("matching to %q", testcase.in), func(t *testing.T) {
			got := framework.MatchImportPath(testcase.in)
			if got != testcase.want {
				t.Errorf("want %t, got %t", testcase.want, got)
			}
		})
	}
}

func Test_NodeConverter_ToEndpoint(t *testing.T) {
	framework := goji.New()

	for _, testcase := range []struct {
		name string
		in   string
		want []endpoint.Endpoint
	}{
		{
			name: "Static Paths",
			in:   `pat.Get("/users")`,
			want: []endpoint.Endpoint{endpoint.New("GET", "/users", "/users")},
		},
		{
			name: "Named Matches",
			in:   `pat.Delete("/users/:user_id")`,
			want: []endpoint.Endpoint{endpoint.New("DELETE", "/users/:user_id", "/users/:user_id")},
		},
		{
			name: "Prefix Matches",
			in:   `pat.Post("/users/files/*")`,
			want: []endpoint.Endpoint{endpoint.New("POST", "/users/files/*", "/users/files/*")},
		},
		{
			name: "Non API Endpoint Register",
			in:   `fmt.Println("/users")`,
			want: nil,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(testcase.in)
			if err != nil {
				t.Error(err)
				return
			}

			got := framework.Extract(expr)
			if !reflect.DeepEqual(got, testcase.want) {
				t.Errorf("want %q, got %q", testcase.want, got)
			}
		})
	}
}
