package parser

import (
	"reflect"
	"testing"
)

func TestParsedRoutes_AddRules(t *testing.T) {
	type fields struct {
		routes map[string]map[string][]Resp
	}
	type args struct {
		ar     Arg
		bodies []Resp
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		afterCall ParsedRoutes
		wantErr   bool
	}{
		{
			name: "Add route to empty set of fields",
			fields: fields{
				routes: map[string]map[string][]Resp{},
			},
			args: args{
				ar: Arg{
					reqType: "GET",
					path:    "/",
				},
				bodies: []Resp{Resp{}},
			},
			afterCall: ParsedRoutes{
				Routes: map[string]map[string][]Resp{
					"/": {
						"GET": []Resp{Resp{}},
					},
				},
			},
		},
		{
			name: "Add Rule Twice",
			fields: fields{
				routes: map[string]map[string][]Resp{
					"/": {
						"GET": []Resp{Resp{}},
					},
				},
			},
			args: args{
				ar: Arg{
					reqType: "GET",
					path:    "/",
				},
				bodies: []Resp{Resp{}},
			},
			afterCall: ParsedRoutes{
				Routes: map[string]map[string][]Resp{
					"/": {
						"GET": []Resp{Resp{}},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rts := &ParsedRoutes{
				Routes: tt.fields.routes,
			}
			err := rts.AddRules(tt.args.ar, tt.args.bodies)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseBurritoFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(rts, &tt.afterCall) {
				t.Errorf("AddRules = %v, want %v", rts, &tt.afterCall)
			}
		})
	}
}

func TestParseBurritoFile(t *testing.T) {
	type args struct {
		filepath string
	}
	tests := []struct {
		name    string
		args    args
		want    ParsedRoutes
		wantErr bool
	}{
		{
			name: "Simple Test",
			args: args{
				filepath: "../../test_burr/single_line.burr",
			},
			want: ParsedRoutes{
				Routes: map[string]map[string][]Resp{
					"/": {
						"GET": []Resp{
							Resp{
								RespType: "FILE",
								Body:     "hello",
							},
						},
					},
				},
			},
		},
		{
			name: "Two-line Test",
			args: args{
				filepath: "../../test_burr/two_lines.burr",
			},
			want: ParsedRoutes{
				Routes: map[string]map[string][]Resp{
					"/": {
						"GET": []Resp{
							Resp{
								RespType: "FILE",
								Body:     "hello",
							},
						},
					},
					"/hello": {
						"GET": []Resp{
							Resp{
								RespType: "FILE",
								Body:     "hello",
							},
						},
					},
				},
			},
		},
		{
			name: "Test with comments",
			args: args{
				filepath: "../../test_burr/include_comments.burr",
			},
			want: ParsedRoutes{
				Routes: map[string]map[string][]Resp{
					"/": {
						"GET": []Resp{
							Resp{
								RespType: "FILE",
								Body:     "hello",
							},
						},
					},
					"/hello": {
						"GET": []Resp{
							Resp{
								RespType: "FILE",
								Body:     "hello",
							},
						},
						"PUT": []Resp{
							Resp{
								RespType: "FILE",
								Body:     "hello",
							},
						},
					},
					"/zesty": {
						"GET": []Resp{
							Resp{
								RespType: "JSON",
								Body:     map[string]interface{}{
									"breakfast" : "burrito",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Test invalid line",
			args: args{
				filepath: "../../test_burr/invalid_line.burr",
			},
			want:    ParsedRoutes{},
			wantErr: true,
		},
		{
			name: "Test invalid line",
			args: args{
				filepath: "../../test_burr/syntax_error_in_args.burr",
			},
			want:    ParsedRoutes{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseBurritoFile(tt.args.filepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseBurritoFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseBurritoFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
