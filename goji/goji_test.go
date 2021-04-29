package goji_test

import (
	"fmt"
	"go/parser"
	"reflect"
	"testing"

	"github.com/progfay/quden/goji"
	"github.com/progfay/quden/util"
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
		want []util.Endpoint
	}{
		{
			name: "Static Paths",
			in:   `pat.Get("/users")`,
			want: []util.Endpoint{util.NewEndpoint("GET /users", "GET /users")},
		},
		{
			name: "Named Matches",
			in:   `pat.Delete("/users/:user_id")`,
			want: []util.Endpoint{util.NewEndpoint("DELETE /users/:user_id", "DELETE /users/:user_id")},
		},
		{
			name: "Prefix Matches",
			in:   `pat.Post("/users/files/*")`,
			want: []util.Endpoint{util.NewEndpoint("POST /users/files/*", "POST /users/files/*")},
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
