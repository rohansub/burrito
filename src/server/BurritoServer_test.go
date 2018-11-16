package server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBurritoServer_Run(t *testing.T) {
	type fields struct {
		Routes *ParsedRoutes
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
		{
			name: "Standard Server Test",
			fields: fields{
				Routes: &ParsedRoutes{
					routes: map[string]map[string][]Resp{
						"/": {
							"GET": []Resp{
								Resp{
									respType: "FILE",
									body:     "../../test_html/hello.html",
								},
							},
						},
						"/hello": {
							"GET": []Resp{
								Resp{
									respType: "FILE",
									body:     "../../test_html/world.html",
								},
							},
							"PUT": []Resp{
								Resp{
									respType: "STR",
									body:     "I am zesty",
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
			NewBurritoServer(tt.fields.Routes)

			// source: https://stackoverflow.com/questions/16154999/how-to-test-http-calls-in-go-using-httptest
			resp := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/", nil)

			http.DefaultServeMux.ServeHTTP(resp, req)
			if p, err := ioutil.ReadAll(resp.Body); err != nil {
				t.Fail()
			} else {
				if !strings.Contains(string(p), "<p>Hello</p>") {
					t.Errorf("header response doen't match:\n%s", p)
				}
			}

			// Test /hello
			req, _ = http.NewRequest("GET", "/hello", nil)

			http.DefaultServeMux.ServeHTTP(resp, req)
			if p, err := ioutil.ReadAll(resp.Body); err != nil {
				t.Fail()
			} else {
				if !strings.Contains(string(p), "<p>World</p>") {
					t.Errorf("header response doen't match:\n%s", p)
				}
			}

			// Test /hello, PUT
			req, _ = http.NewRequest("PUT", "/hello", nil)

			http.DefaultServeMux.ServeHTTP(resp, req)
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
