package gorilla_test

import (
	"fmt"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
	"testing"

	"github.com/progfay/quden/gorilla"
	"github.com/progfay/quden/util"
)

const template = `package main

import "github.com/gorilla/mux"

func main() {
	r := mux.NewRouter()
	%s
	log.Fatal(http.ListenAndServe(":8000", r))
}
`

func Test_MatchImportPath(t *testing.T) {
	framework := gorilla.New()

	for _, testcase := range []struct {
		in   string
		want bool
	}{
		{
			in:   "github.com/gorilla/mux",
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
	framework := gorilla.New()

	for _, testcase := range []struct {
		name string
		in   []string
		want []util.Endpoint
	}{
		{
			name: "Static Paths",
			in:   []string{`r.HandleFunc("/users", handler).Methods("GET")`},
			want: []util.Endpoint{util.NewEndpoint("GET /users", `^GET /users\b`)},
		},
		{
			name: "Variable Paths",
			in:   []string{`r.HandleFunc("/users/{user_id}", handler).Methods("DELETE")`},
			want: []util.Endpoint{util.NewEndpoint("DELETE /users/{user_id}", `^DELETE /users/[^/]+\b`)},
		},
		{
			name: "Non API Endpoint Register",
			in:   []string{`fmt.Println("/users")`},
			want: nil,
		},
		{
			name: "Subrouter",
			in: []string{
				`l := r.PathPrefix("/users").Subrouter()`,
				`l.Methods("GET").HandlerFunc(handler)`,
			},
			want: []util.Endpoint{util.NewEndpoint("GET /users", `^GET /users\B*\b`)},
		},
		{
			name: "Multiple methods",
			in:   []string{`r.HandleFunc("/users", handler).Methods("GET", "POST")`},
			want: []util.Endpoint{
				util.NewEndpoint("GET /users", `^GET /users\b`),
				util.NewEndpoint("POST /users", `^POST /users\b`),
			},
		},
		{
			name: "No specified methods",
			in:   []string{`r.HandleFunc("/users", handler)`},
			want: []util.Endpoint{util.NewEndpoint("/users", `^[^ ]+ /users\b`)},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			source := fmt.Sprintf(template, strings.Join(testcase.in, "\n  "))
			expr, err := parser.ParseFile(token.NewFileSet(), "", source, parser.Mode(0))
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
