package handler

import (
	"reflect"
	"testing"
)

func TestNewZester(t *testing.T) {
	tests := []struct {
		name string
		want *Zester
	}{
		{
			name: "Creation test",
			want: &Zester{
				paths: make([]*PathObject,0),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewZester(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewZester() = %v, want %v", got, tt.want)
			}
		})
	}
}


func TestZester_CheckPaths(t *testing.T) {
	type fields struct {
		paths []*PathObject
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Test simple strings, no conflicts",
			fields: fields {
				paths: []*PathObject{
					{
						str: "/hello/world",
						parts: []*PathSegment {
							NewPathSegment("hello"),
							NewPathSegment("world"),
						},
					},
					{
						str: "/zesty/burrito",
						parts: []*PathSegment {
							NewPathSegment("zesty"),
							NewPathSegment("burrito"),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Test simple strings, one segment matches, no conflict",
			fields: fields {
				paths: []*PathObject{
					{
						str: "/zesty/world",
						parts: []*PathSegment {
							NewPathSegment("zesty"),
							NewPathSegment("world"),
						},
					},
					{
						str: "/zesty/burrito",
						parts: []*PathSegment {
							NewPathSegment("zesty"),
							NewPathSegment("burrito"),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Test simple strings, conflict between equivalent urls",
			fields: fields {
				paths: []*PathObject{
					{
						str: "/zesty/burrito",
						parts: []*PathSegment {
							NewPathSegment("zesty"),
							NewPathSegment("burrito"),
						},
					},
					{
						str: "/zesty/burrito",
						parts: []*PathSegment {
							NewPathSegment("zesty"),
							NewPathSegment("burrito"),
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Test variable conflict with string",
			fields: fields {
				paths: []*PathObject{
					{
						str: "/zesty/:chocolate",
						parts: []*PathSegment {
							NewPathSegment("zesty"),
							NewPathSegment(":chocolate"),
						},
					},
					{
						str: "/zesty/burrito",
						parts: []*PathSegment {
							NewPathSegment("zesty"),
							NewPathSegment("burrito"),
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Test string variable conflict with string",
			fields: fields {
				paths: []*PathObject{
					{
						str: "/zesty/:chocolate:str",
						parts: []*PathSegment {
							NewPathSegment("zesty"),
							NewPathSegment(":chocolate:str"),
						},
					},
					{
						str: "/zesty/burrito",
						parts: []*PathSegment {
							NewPathSegment("zesty"),
							NewPathSegment("burrito"),
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Test int variable no conflict with string",
			fields: fields {
				paths: []*PathObject{
					{
						str: "/zesty/:chocolate:int",
						parts: []*PathSegment {
							NewPathSegment("zesty"),
							NewPathSegment(":chocolate:int"),
						},
					},
					{
						str: "/zesty/burrito",
						parts: []*PathSegment {
							NewPathSegment("zesty"),
							NewPathSegment("burrito"),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Test float variable no conflict with string",
			fields: fields {
				paths: []*PathObject{
					{
						str: "/zesty/:chocolate:float",
						parts: []*PathSegment {
							NewPathSegment("zesty"),
							NewPathSegment(":chocolate:float"),
						},
					},
					{
						str: "/zesty/burrito",
						parts: []*PathSegment {
							NewPathSegment("zesty"),
							NewPathSegment("burrito"),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Test variable-variable conflict",
			fields: fields {
				paths: []*PathObject{
					{
						str: "/zesty/hello/:chocolate:float",
						parts: []*PathSegment {
							NewPathSegment("zesty"),
							NewPathSegment("hello"),
							NewPathSegment(":chocolate:float"),
						},
					},
					{
						str: "/zesty/hello/:burrito",
						parts: []*PathSegment {
							NewPathSegment("zesty"),
							NewPathSegment("hello"),
							NewPathSegment(":burrito"),
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Test different lengths, no conflict",
			fields: fields {
				paths: []*PathObject{
					{
						str: "/zesty/hello/:chocolate:float",
						parts: []*PathSegment {
							NewPathSegment("zesty"),
							NewPathSegment("hello"),
							NewPathSegment(":chocolate:float"),
						},
					},
					{
						str: "/zesty/hello",
						parts: []*PathSegment {
							NewPathSegment("zesty"),
							NewPathSegment("hello"),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Test many paths, no conflict",
			fields: fields {
				paths: []*PathObject{
					{
						str: "/zesty/hello/:chocolate:float",
						parts: []*PathSegment {
							NewPathSegment("zesty"),
							NewPathSegment("hello"),
							NewPathSegment(":chocolate:float"),
						},
					},
					{
						str: "/zesty/hello",
						parts: []*PathSegment {
							NewPathSegment("zesty"),
							NewPathSegment("hello"),
						},
					},
					{
						str: "/zesty",
						parts: []*PathSegment {
							NewPathSegment("zesty"),
						},
					},
					{
						str: "/cookie/:dough",
						parts: []*PathSegment {
							NewPathSegment("cookie"),
							NewPathSegment(":dough"),

						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Test many paths, conflict",
			fields: fields {
				paths: []*PathObject{
					{
						str: "/zesty/hello/:chocolate:float",
						parts: []*PathSegment {
							NewPathSegment("zesty"),
							NewPathSegment("hello"),
							NewPathSegment(":chocolate:float"),
						},
					},
					{
						str: "/zesty/hello",
						parts: []*PathSegment {
							NewPathSegment("zesty"),
							NewPathSegment("hello"),
						},
					},
					{
						str: "/zesty",
						parts: []*PathSegment {
							NewPathSegment("zesty"),
						},
					},
					{
						str: "/cookie/:dough",
						parts: []*PathSegment {
							NewPathSegment("cookie"),
							NewPathSegment(":dough"),

						},
					},
					{
						str: "/:cookie/:dough",
						parts: []*PathSegment {
							NewPathSegment(":cookie"),
							NewPathSegment(":dough"),
						},
					},
				},
			},
			wantErr: true,
		},


	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := &Zester{
				paths: tt.fields.paths,
			}
			if err := z.CheckPaths(); (err != nil) != tt.wantErr {
				t.Errorf("Zester.CheckPaths() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_isConflicting(t *testing.T) {
	type args struct {
		one *PathObject
		two *PathObject
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isConflicting(tt.args.one, tt.args.two); got != tt.want {
				t.Errorf("isConflicting() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getStringMatching(t *testing.T) {
	type args struct {
		one *PathSegment
		two *PathSegment
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test two Non variables same strings",
			args: args {
				one: NewPathSegment("hello"),
				two: NewPathSegment("hello"),
			},
			want: true,
		},
		{
			name: "Test two Non variables different strings",
			args: args {
				one: NewPathSegment("hello"),
				two: NewPathSegment("world"),
			},
			want: false,
		},
		{
			name: "Test non variable compared with variable",
			args: args {
				one: NewPathSegment("hello"),
				two: NewPathSegment(":world"),
			},
			want: true,
		},
		{
			name: "Test non variable compared with explicit string variable",
			args: args {
				one: NewPathSegment("hello"),
				two: NewPathSegment(":world:str"),
			},
			want: true,
		},
		{
			name: "Test non variable compared with explicit int variable",
			args: args {
				one: NewPathSegment("hello"),
				two: NewPathSegment(":world:int"),
			},
			want: false,
		},
		{
			name: "Test non variable compared with explicit int variable",
			args: args {
				one: NewPathSegment("hello"),
				two: NewPathSegment(":world:int"),
			},
			want: false,
		},
		{
			name: "Test non variable compared with explicit float variable",
			args: args {
				one: NewPathSegment("hello"),
				two: NewPathSegment(":world:float"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSegmentMatch(tt.args.one, tt.args.two); got != tt.want {
				t.Errorf("getSegmentMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}


