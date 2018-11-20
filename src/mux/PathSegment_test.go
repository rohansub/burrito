package mux

import (
	"reflect"
	"testing"
)

func TestNewPathSegment(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want *PathSegment
	}{
		{
			name: "Test No vars",
			args: args{
				"hello",
			},
			want: &PathSegment{
				mustMatch: true,
				segStr:    "hello",
				typeMatch: "",
			},
		},
		{
			name: "Test nothing in string",
			args: args{
				"",
			},
			want: &PathSegment{
				mustMatch: true,
				segStr:    "",
				typeMatch: "",
			},
		},
		{
			name: "Test var in string",
			args: args{
				":hello",
			},
			want: &PathSegment{
				mustMatch: false,
				varName:   "hello",
			},
		},
		{
			name: "Test var in string with type",
			args: args{
				":hello:int",
			},
			want: &PathSegment{
				mustMatch: false,
				varName:   "hello",
				typeMatch: "int",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPathSegment(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPathSegment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPathSegment_SegMatchAndExtractVars(t *testing.T) {
	type fields struct {
		mustMatch bool
		segStr    string
		typeMatch string
		varName   string
	}
	type args struct {
		str string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  *EnvEntry
	}{
		{
			name: "Test normal string, no vars, match success",
			fields: fields{
				mustMatch: true,
				segStr:    "hello",
			},
			args: args{
				str: "hello",

			},
			want: true,
			want1: nil,
		},
		{
			name: "Test normal string, no vars, not matching",
			fields: fields{
				mustMatch: true,
				segStr:    "hello",
			},
			args: args{
				str: "f",

			},
			want: false,
			want1: nil,
		},
		{
			name: "Test var string, no type, matching",
			fields: fields{
				mustMatch: false,
				varName:    "burrito",
			},
			args: args{
				str: "f",

			},
			want: true,
			want1: &EnvEntry{
				name: "burrito",
				isStr: true,
				valStr: "f",
			},
		},
		{
			name: "Test var with int type, matching",
			fields: fields{
				mustMatch: false,
				varName:    "burrito",
				typeMatch: "int",
			},
			args: args{
				str: "44",

			},
			want: true,
			want1: &EnvEntry{
				name: "burrito",
				isInt: true,
				valInt: 44,
			},
		},
		{
			name: "Test var with int type, not matching",
			fields: fields{
				mustMatch: false,
				varName:    "burrito",
				typeMatch: "int",
			},
			args: args{
				str: "notstring",

			},
			want: false,
			want1: nil,
		},
		{
			name: "Test var with float type, matching",
			fields: fields{
				mustMatch: false,
				varName:   "burrito",
				typeMatch: "float",
			},
			args: args{
				str: "44.3",
			},
			want: true,
			want1: &EnvEntry{
				name: "burrito",
				isFlt: true,
				valFlt: 44.3,
			},
		},
		{
			name: "Test var with float type, not matching",
			fields: fields{
				mustMatch: false,
				varName:    "burrito",
				typeMatch: "float",
			},
			args: args{
				str: "notstring",
			},
			want: false,
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &PathSegment{
				mustMatch: tt.fields.mustMatch,
				segStr:    tt.fields.segStr,
				typeMatch: tt.fields.typeMatch,
				varName:   tt.fields.varName,
			}
			got, got1 := ps.SegMatchAndExtractVars(tt.args.str)
			if got != tt.want {
				t.Errorf("PathSegment.SegMatchAndExtractVars() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PathSegment.SegMatchAndExtractVars() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
