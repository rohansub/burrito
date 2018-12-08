package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/rcsubra2/burrito/src/db"
	"github.com/rcsubra2/burrito/src/parser"
	"github.com/rcsubra2/burrito/src/utils"
)

func TestBurritoServer_Run(t *testing.T) {
	type fields struct {
		Routes   *parser.ParsedRoutes
		MockInit map[string]string
	}
	type arg struct {
		method   string
		uri      string
		wantType string
		want     interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   []arg
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
					method:   "GET",
					uri:      "/",
					want:     "<p>Hello</p>",
					wantType: "text/html",
				},
				{
					method:   "GET",
					uri:      "/hello",
					want:     "<p>World</p>",
					wantType: "text/html",
				},
				{
					method:   "PUT",
					uri:      "/hello",
					want:     "I am zesty",
					wantType: "text/html",
				},
			},
		},
		{
			name: "Serve Template data, given url variables",
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
					method:   "GET",
					uri:      "/",
					want:     "<p>Hello</p>",
					wantType: "text/html",
				},
				{
					method:   "GET",
					uri:      "/zesty/fool",
					want:     "<h>fool</h>",
					wantType: "text/html",
				},
			},
		},
		{
			name: "GET JSON data",
			fields: fields{
				Routes: &parser.ParsedRoutes{
					Routes: map[string]map[string][]parser.Resp{
						"/": {
							"GET": []parser.Resp{
								parser.Resp{
									RespType: "JSON",
									Body: map[string]interface{}{
										"zesty": "burrito",
									},
								},
							},
						},
					},
				},
			},
			args: []arg{
				{
					method:   "GET",
					uri:      "/",
					wantType: "app/json",
					want: map[string]interface{}{
						"zesty": "burrito",
					},
				},
			},
		},
		{
			name: "GET Redis data",
			fields: fields{
				Routes: &parser.ParsedRoutes{
					Routes: map[string]map[string][]parser.Resp{
						"/": {
							"GET": []parser.Resp{
								{
									RespType: "DB",
									DBReq: &db.GetReq{
										ArgNames: []db.Param{
											{
												IsString: true,
												Val:      "zesty",
											},
										},
									},
								},
							},
						},
					},
				},
				MockInit: map[string]string{
					"zesty": "burrito",
				},
			},
			args: []arg{
				{
					method:   "GET",
					uri:      "/",
					wantType: "app/json",
					want: map[string]interface{}{
						"zesty": "burrito",
					},
				},
			},
		},
		{
			name: "GET Redis data, url variable test",
			fields: fields{
				Routes: &parser.ParsedRoutes{
					Routes: map[string]map[string][]parser.Resp{
						"/:zesty": {
							"GET": []parser.Resp{
								{
									RespType: "DB",
									DBReq: &db.GetReq{
										ArgNames: []db.Param{
											{
												IsString: false,
												Val:      "zesty",
											},
										},
									},
								},
							},
						},
					},
				},
				MockInit: map[string]string{
					"hello": "burrito",
				},
			},
			args: []arg{
				{
					method:   "GET",
					uri:      "/hello",
					wantType: "app/json",
					want: map[string]interface{}{
						"hello": "burrito",
					},
				},
			},
		},
		{
			name: "GET Redis data, multiple items",
			fields: fields{
				Routes: &parser.ParsedRoutes{
					Routes: map[string]map[string][]parser.Resp{
						"/:zesty": {
							"GET": []parser.Resp{
								{
									RespType: "DB",
									DBReq: &db.GetReq{
										ArgNames: []db.Param{
											{
												IsString: false,
												Val:      "zesty",
											},
											{
												IsString: true,
												Val:      "quesadilla",
											},
										},
									},
								},
							},
						},
					},
				},
				MockInit: map[string]string{
					"hello":      "burrito",
					"quesadilla": "cheese",
				},
			},
			args: []arg{
				{
					method:   "GET",
					uri:      "/hello",
					wantType: "app/json",
					want: map[string]interface{}{
						"hello":      "burrito",
						"quesadilla": "cheese",
					},
				},
			},
		},
		{
			name: "GET Redis data, multiple items, not all exist",
			fields: fields{
				Routes: &parser.ParsedRoutes{
					Routes: map[string]map[string][]parser.Resp{
						"/:zesty": {
							"GET": []parser.Resp{
								{
									RespType: "DB",
									DBReq: &db.GetReq{
										ArgNames: []db.Param{
											{
												IsString: false,
												Val:      "zesty",
											},
											{
												IsString: true,
												Val:      "quesadilla",
											},
										},
									},
								},
							},
						},
					},
				},
				MockInit: map[string]string{
					"hello": "burrito",
				},
			},
			args: []arg{
				{
					method:   "GET",
					uri:      "/hello",
					wantType: "app/json",
					want: map[string]interface{}{
						"hello": "burrito",
					},
				},
			},
		},
		{
			name: "SET Redis data, single item",
			fields: fields{
				Routes: &parser.ParsedRoutes{
					Routes: map[string]map[string][]parser.Resp{
						"/:zesty": {
							"PUT": []parser.Resp{
								{
									RespType: "DB",
									DBReq: &db.SetReq{
										ArgNames: []utils.Pair{
											{
												Fst: db.Param{
													IsString: false,
													Val:      "zesty",
												},
												Snd: db.Param{
													IsString: true,
													Val:      "quesadilla",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				MockInit: map[string]string{
					"hello": "burrito",
				},
			},
			args: []arg{
				{
					method:   "PUT",
					uri:      "/hello",
					wantType: "app/json",
					want:     map[string]interface{}{},
				},
			},
		},
		{
			name: "SET chained with Get Redis data, multiple items",
			fields: fields{
				Routes: &parser.ParsedRoutes{
					Routes: map[string]map[string][]parser.Resp{
						"/:zesty": {
							"PUT": []parser.Resp{
								{
									RespType: "DB",
									DBReq: &db.SetReq{
										ArgNames: []utils.Pair{
											{
												Fst: db.Param{
													IsString: false,
													Val:      "zesty",
												},
												Snd: db.Param{
													IsString: true,
													Val:      "quesadilla",
												},
											},
											{
												Fst: db.Param{
													IsString: true,
													Val:      "burrito",
												},
												Snd: db.Param{
													IsString: true,
													Val:      "supreme",
												},
											},
										},
									},
								},
								{
									RespType: "DB",
									DBReq: &db.GetReq{
										ArgNames: []db.Param{
											{
												IsString: false,
												Val:      "zesty",
											},
											{
												IsString: true,
												Val:      "burrito",
											},
										},
									},
								},
							},
						},
					},
				},
				MockInit: map[string]string{
					"hello":   "different",
					"burrito": "notsame",
				},
			},
			args: []arg{
				{
					method:   "PUT",
					uri:      "/hello",
					wantType: "app/json",
					want: map[string]interface{}{
						"hello":   "quesadilla",
						"burrito": "supreme",
					},
				},
			},
		},
		{
			name: "GET Redis, chained with SET Redis, chained with html",
			fields: fields{
				Routes: &parser.ParsedRoutes{
					Routes: map[string]map[string][]parser.Resp{
						"/:name": {
							"PUT": []parser.Resp{
								{
									RespType: "DB",
									DBReq: &db.SetReq{
										ArgNames: []utils.Pair{
											{
												Fst: db.Param{
													IsString: true,
													Val:      "name",
												},
												Snd: db.Param{
													IsString: false,
													Val:      "name",
												},
											},
											{
												Fst: db.Param{
													IsString: true,
													Val:      "greeting",
												},
												Snd: db.Param{
													IsString: true,
													Val:      "hello",
												},
											},
										},
									},
								},
								{
									RespType: "DB",
									DBReq: &db.GetReq{
										ArgNames: []db.Param{
											{
												IsString: true,
												Val:      "name",
											},
											{
												IsString: true,
												Val:      "greeting",
											},
										},
									},
								},
								{
									RespType: "FILE",
									Body:     "../../test_html/user.html",
								},
							},
						},
					},
				},
				MockInit: map[string]string{},
			},
			args: []arg{
				{
					method:   "PUT",
					uri:      "/rohan",
					wantType: "text/html",
					want:     "<h>rohan</h>\n<h>rohan</h>\n<h>hello</h>",
				},
			},
		},
		{
			name: "Serve Template data, given form data",
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
						"/zesty": {
							"GET": []parser.Resp{
								parser.Resp{
									RespType: "FILE",
									Body:     "../../test_html/template2.html",
								},
							},
						},
					},
				},
			},
			args: []arg{
				{
					method:   "GET",
					uri:      "/",
					want:     "<p>Hello</p>",
					wantType: "text/html",
				},
				{
					method:   "GET",
					uri:      "/zesty?burrito=fool",
					want:     "<h>fool</h>",
					wantType: "text/html",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := NewBurritoServer(tt.fields.Routes, tt.fields.MockInit)

			// source: https://stackoverflow.com/questions/16154999/how-to-test-http-calls-in-go-using-httptest
			resp := httptest.NewRecorder()

			for _, ar := range tt.args {
				req, _ := http.NewRequest(ar.method, ar.uri, nil)

				b.router.ServeHTTP(resp, req)
				cType := resp.Header().Get("Content-type")

				if !reflect.DeepEqual(resp.Header().Get("Content-type"), ar.wantType) {
					t.Errorf("Wanted Content-type %s - Got %s", ar.wantType, cType)
				}
				if p, err := ioutil.ReadAll(resp.Body); err != nil {
					t.Fail()
				} else {
					if ar.wantType == "text/html" && !reflect.DeepEqual(string(p), ar.want) {
						t.Errorf("Wanted %s - Got %s", ar.want, p)
					} else if ar.wantType == "app/json" {
						var d interface{}
						err := json.Unmarshal([]byte(p), &d)
						if err != nil {
							t.Errorf("JSON data invalid: Wanted %s - Got %s", ar.want, p)
						} else if !reflect.DeepEqual(ar.want, d) {
							t.Errorf("Wanted %s - Got %s", ar.want, d)
						}
					}
				}
			}
		})

	}
}

func TestNewBurritoServer(t *testing.T) {
	type args struct {
		rts      *parser.ParsedRoutes
		mockData map[string]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test Create, with no route conflicts",
			args: args {
				rts: &parser.ParsedRoutes{
					Routes: map[string]map[string][]parser.Resp{
						"/": {
							"GET": []parser.Resp{},
						},
						"/hello": {
							"GET": []parser.Resp{},
						},
						"/hello/:world": {
							"GET": []parser.Resp{},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Test Create, with route conflicts",
			args: args {
				rts: &parser.ParsedRoutes{
					Routes: map[string]map[string][]parser.Resp{
						"/": {
							"GET": []parser.Resp{},
						},
						"/hello": {
							"GET": []parser.Resp{},
						},
						"/hello/:world": {
							"GET": []parser.Resp{},
						},
						"/:chicken/nugget": {
							"GET": []parser.Resp{},
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewBurritoServer(tt.args.rts, tt.args.mockData)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBurritoServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
