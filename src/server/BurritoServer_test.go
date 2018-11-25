package server

import (
	"github.com/rcsubra2/burrito/src/parser"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBurritoServer_Run(t *testing.T) {
	type fields struct {
		Routes *parser.ParsedRoutes
	}
	type arg struct {
		method string
		uri    string
		want   string
	}
	tests := []struct {
		name   string
		fields fields
		args []arg
	}{
		{
			name: "Standard Server String and file data only",
			fields: fields{
				Routes: &parser.ParsedRoutes{
					Routes: map[string]map[string][]parser.Resp{
						"/": {
							"GET": []parser.Resp{
								parser.Resp{
									RespType: "FILE",
									Body:     "../../test_html/hello.html",
								},
							},
						},
						"/hello": {
							"GET": []parser.Resp{
								parser.Resp{
									RespType: "FILE",
									Body:     "../../test_html/world.html",
								},
							},
							"PUT": []parser.Resp{
								parser.Resp{
									RespType: "STR",
									Body:     "I am zesty",
								},
							},
						},
					},
				},
			},
			args: []arg{
				{
					method: "GET",
					uri: "/",
					want: "<p>Hello</p>",
				},
				{
					method: "GET",
					uri: "/hello",
					want: "<p>World</p>",
				},
				{
					method: "PUT",
					uri: "/hello",
					want: "I am zesty",
				},
			},
		},
		{
			name: "Standard Server String data Generics",
			fields: fields{
				Routes: &parser.ParsedRoutes{
					Routes: map[string]map[string][]parser.Resp{
						"/": {
							"GET": []parser.Resp{
								parser.Resp{
									RespType: "FILE",
									Body:     "../../test_html/hello.html",
								},
							},
						},
						"/zesty/:burrito": {
							"GET": []parser.Resp{
								parser.Resp{
									RespType: "FILE",
									Body:     "../../test_html/template.html",
								},
							},
						},
					},
				},
			},
			args: []arg{
				{
					method: "GET",
					uri: "/",
					want: "<p>Hello</p>",
				},
				{
					method: "GET",
					uri: "/zesty/fool",
					want: "<h>fool</h>",
				},
			},

		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBurritoServer(tt.fields.Routes)

			// source: https://stackoverflow.com/questions/16154999/how-to-test-http-calls-in-go-using-httptest
			resp := httptest.NewRecorder()

			for _, ar := range tt.args {
				req, _ := http.NewRequest(ar.method, ar.uri, nil)

				b.router.ServeHTTP(resp, req)
				if p, err := ioutil.ReadAll(resp.Body); err != nil {
					t.Fail()
				} else {
					if string(p) != ar.want {
						t.Errorf("Want %s - Got %s", ar.want, p)
					}
				}
			}



		})
	}
}
