package echo_test

import (
	"fmt"
	"go/parser"
	"reflect"
	"testing"

	"github.com/progfay/quden/echo"
	"github.com/progfay/quden/endpoint"
)

func Test_MatchImportPath(t *testing.T) {
	framework := echo.New()

	for _, testcase := range []struct {
		in   string
		want bool
	}{
		{
			in:   "github.com/labstack/echo/v4",
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
	framework := echo.New()

	for _, testcase := range []struct {
		name string
		in   string
		want *endpoint.Endpoint
	}{
		{
			name: "Static",
			in:   `e.GET("/users", handler)`,
			want: endpoint.New("GET", "/users"),
		},
		{
			name: "Param",
			in:   `e.DELETE("/users/:user_id", handler)`,
			want: endpoint.New("DELETE", "/users/:user_id"),
		},
		{
			name: "Match Any",
			in:   `e.POST("/users/files/*", handler)`,
			want: endpoint.New("POST", "/users/files/*"),
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

			got := framework.NewNodeConverter().ToEndpoint(expr)
			if !reflect.DeepEqual(got, testcase.want) {
				t.Errorf("want %q, got %q", testcase.want, got)
			}
		})
	}
}
