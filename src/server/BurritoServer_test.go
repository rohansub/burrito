package server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/rcsubra2/burrito/src/parser"
)

func TestBurritoServer_Run(t *testing.T) {
	type fields struct {
		Routes *parser.ParsedRoutes
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
		{
			name: "Standard Server Test",
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBurritoServer(tt.fields.Routes)

			// source: https://stackoverflow.com/questions/16154999/how-to-test-http-calls-in-go-using-httptest
			resp := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/", nil)

			b.router.ServeHTTP(resp, req)
			if p, err := ioutil.ReadAll(resp.Body); err != nil {
				t.Fail()
			} else {
				if !strings.Contains(string(p), "<p>Hello</p>") {
					t.Errorf("header response doen't match:\n%s", p)
				}
			}

			// Test /hello
			req, _ = http.NewRequest("GET", "/hello", nil)

			b.router.ServeHTTP(resp, req)
			if p, err := ioutil.ReadAll(resp.Body); err != nil {
				t.Fail()
			} else {
				if !strings.Contains(string(p), "<p>World</p>") {
					t.Errorf("header response doen't match:\n%s", p)
				}
			}

			// Test /hello, PUT
			req, _ = http.NewRequest("PUT", "/hello", nil)

			b.router.ServeHTTP(resp, req)
			if p, err := ioutil.ReadAll(resp.Body); err != nil {
				t.Fail()
			} else {
				if !strings.Contains(string(p), "I am zesty") {
					t.Errorf("header response doen't match:\n%s", p)
				}
			}

		})
	}
}
