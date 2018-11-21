package handler

import (
	"reflect"
	"testing"

	"github.com/rcsubra2/burrito/src/environment"
)

func Test_route_Match(t *testing.T) {
	type fields struct {
		pattern []*PathSegment
		handler BurritoHandler
	}
	type args struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  *environment.Env
	}{
		{
			name: "Route no variables, route matches",
			fields: fields {
				pattern: []*PathSegment{
					{
						mustMatch: true,
						segStr: "zesty",
					},
					{
						mustMatch: true,
						segStr: "burrito",
					},
				},
			},
			args: args {
				path: "/zesty/burrito",
			},
			want: true,
			want1: func() (*environment.Env){
				e := environment.CreateEnv()
				return e
			}(),
		},
		{
			name: "Route no variables, route doesn't match",
			fields: fields {
				pattern: []*PathSegment{
					{
						mustMatch: true,
						segStr: "zesty",
					},
					{
						mustMatch: true,
						segStr: "burrito",
					},
				},
			},
			args: args {
				path: "/zesty/notburrito",
			},
			want: false,
			want1: nil,
		},
		{
			name: "Route no variables, route doesn't match due to length",
			fields: fields {
				pattern: []*PathSegment{
					{
						mustMatch: true,
						segStr: "zesty",
					},
					{
						mustMatch: true,
						segStr: "burrito",
					},
				},
			},
			args: args {
				path: "/zesty",
			},
			want: false,
			want1: nil,
		},
		{
			name: "Route int variable, route matches",
			fields: fields {
				pattern: []*PathSegment{
					{
						mustMatch: true,
						segStr: "zesty",
					},
					{
						mustMatch: false,
						varName: "burrito",
						typeMatch: "int",
					},
				},
			},
			args: args {
				path: "/zesty/64",
			},
			want: true,
			want1: func() (*environment.Env){
				e := environment.CreateEnv()
				e.Add(*environment.CreateIntEntry("burrito", 64))
				return e
			}(),
		},
		{
			name: "Route float variable, route matches",
			fields: fields {
				pattern: []*PathSegment{
					{
						mustMatch: true,
						segStr: "zesty",
					},
					{
						mustMatch: false,
						varName: "burrito",
						typeMatch: "float",
					},
				},
			},
			args: args {
				path: "/zesty/64.3",
			},
			want: true,
			want1: func() (*environment.Env){
				e := environment.CreateEnv()
				e.Add(*environment.CreateFloatEntry("burrito", 64.3))
				return e
			}(),
		},
		{
			name: "Route string variable, route matches",
			fields: fields {
				pattern: []*PathSegment{
					{
						mustMatch: true,
						segStr: "zesty",
					},
					{
						mustMatch: false,
						varName: "burrito",
						typeMatch: "string",
					},
				},
			},
			args: args {
				path: "/zesty/hello",
			},
			want: true,
			want1: func() (*environment.Env){
				e := environment.CreateEnv()
				e.Add(*environment.CreateStringEntry("burrito", "hello"))
				return e
			}(),
		},
		{
			name: "Route int variable, no route match",
			fields: fields {
				pattern: []*PathSegment{
					{
						mustMatch: true,
						segStr: "zesty",
					},
					{
						mustMatch: false,
						varName: "burrito",
						typeMatch: "int",
					},
				},
			},
			args: args {
				path: "/zesty/hello",
			},
			want: false,
			want1: nil,
		},
		{
			name: "Route float variable, no route match",
			fields: fields {
				pattern: []*PathSegment{
					{
						mustMatch: true,
						segStr: "zesty",
					},
					{
						mustMatch: false,
						varName: "burrito",
						typeMatch: "float",
					},
				},
			},
			args: args {
				path: "/zesty/hello",
			},
			want: false,
			want1: nil,
		},
		{
			name: "Multiple variables,  match",
			fields: fields {
				pattern: []*PathSegment{
					{
						mustMatch: false,
						varName: "zesty",
						typeMatch: "string",
					},
					{
						mustMatch: false,
						varName: "burrito",
						typeMatch: "int",
					},
				},
			},
			args: args {
				path: "/breakfast/42",
			},
			want: true,
			want1: func() (*environment.Env){
				e := environment.CreateEnv()
				e.Add(*environment.CreateStringEntry("zesty", "breakfast"))
				e.Add(*environment.CreateIntEntry("burrito", 42))
				return e
			}(),
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &route{
				pattern: tt.fields.pattern,
				handler: tt.fields.handler,
			}
			got, got1 := r.Match(tt.args.path)
			if got != tt.want {
				t.Errorf("route.Match() got = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("route.Match() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
