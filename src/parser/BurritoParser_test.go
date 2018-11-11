package parser

import (
	"reflect"
	"testing"
)

func TestParsedRoutes_AddRules(t *testing.T) {
	type fields struct {
		routes map[Arg][]Resp
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
	}{
		{
			name: "Add route to empty set of fields",
			fields: fields{
				routes: map[Arg][]Resp{},
			},
			args: args{
				ar: Arg{
					reqType: "GET",
					path:    "/",
				},
				bodies: []Resp{Resp{}},
			},
			afterCall: ParsedRoutes{
				routes: map[Arg][]Resp{
					Arg{path: "/", reqType: "GET"}: []Resp{Resp{}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rts := &ParsedRoutes{
				routes: tt.fields.routes,
			}
			rts.AddRules(tt.args.ar, tt.args.bodies)
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
				filepath: "../../test_docs/single_line.burr",
			},
			want: ParsedRoutes{
				routes: map[Arg][]Resp{
					Arg{reqType: "GET", path: "/"}: []Resp{
						Resp{
							respType: "FILE",
							body:     "hello",
						},
					},
				},
			},
		},
		{
			name: "Two-line Test",
			args: args{
				filepath: "../../test_docs/two_lines.burr",
			},
			want: ParsedRoutes{
				routes: map[Arg][]Resp{
					Arg{reqType: "GET", path: "/"}: []Resp{
						Resp{
							respType: "FILE",
							body:     "hello",
						},
					},
					Arg{reqType: "GET", path: "/hello"}: []Resp{
						Resp{
							respType: "FILE",
							body:     "hello",
						},
					},
				},
			},
		},
		{
			name: "Test with comments",
			args: args{
				filepath: "../../test_docs/include_comments.burr",
			},
			want: ParsedRoutes{
				routes: map[Arg][]Resp{
					Arg{reqType: "GET", path: "/"}: []Resp{
						Resp{
							respType: "FILE",
							body:     "hello",
						},
					},
					Arg{reqType: "GET", path: "/hello"}: []Resp{
						Resp{
							respType: "FILE",
							body:     "hello",
						},
					},
					Arg{reqType: "PUT", path: "/hello"}: []Resp{
						Resp{
							respType: "FILE",
							body:     "hello",
						},
					},
				},
			},
		},
		{
			name: "Test invalid line",
			args: args{
				filepath: "../../test_docs/invalid_line.burr",
			},
			want:    ParsedRoutes{},
			wantErr: true,
		},
		{
			name: "Test invalid line",
			args: args{
				filepath: "../../test_docs/syntax_error_in_args.burr",
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
