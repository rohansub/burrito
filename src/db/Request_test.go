package db

import (
	"reflect"
	"testing"

	"github.com/rcsubra2/burrito/src/environment"
)

func TestParam_GetValue(t *testing.T) {
	env1 := environment.CreateEnv()
	env2 := environment.CreateEnv()

	env1.Add(*environment.CreateStringEntry("hello", "world"))
	env2.Add(*environment.CreateStringEntry("hello", "world2"))
	env2.Add(*environment.CreateStringEntry("hello2", "world2"))



	type fields struct {
		IsString bool
		Val      string
	}
	type args struct {
		envs []*environment.Env
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		wantOk  bool
	}{
		{
			name: "Test String",
			fields: fields {
				IsString: true,
				Val: "hello",
			},
			args: args {
				envs: []*environment.Env{env1, env2},
			},
			want: "hello",
			wantOk: true,
		},
		{
			name: "Test Variable",
			fields: fields {
				IsString: false,
				Val: "hello",
			},
			args: args {
				envs: []*environment.Env{env1, env2},
			},
			want: "world",
			wantOk: true,
		},
		{
			name: "Test Variable Not in entry",
			fields: fields {
				IsString: false,
				Val: "nothere",
			},
			args: args {
				envs: []*environment.Env{env1, env2},
			},
			want: "",
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Param{
				IsString: tt.fields.IsString,
				Val:      tt.fields.Val,
			}
			got, ok := p.GetValue(tt.args.envs)
			if (ok != false) != tt.wantOk  {
				t.Errorf("ok = %v want %v", ok, tt.wantOk)
			}
			if got != tt.want {
				t.Errorf("Param.GetValue() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestCreateDBGetReq(t *testing.T) {
	type args struct {
		argStrs []string
	}
	tests := []struct {
		name string
		args args
		want *GetReq
	}{
		{
			name: "Test Get with variables and strings",
			args: args{
				argStrs: []string{"'hello'","variable", "var2", "'string'"},
			},
			want: &GetReq{
				ArgNames: []Param{
					{
						IsString: true,
						Val: "hello",
					},
					{
						IsString: false,
						Val: "variable",
					},
					{
						IsString: false,
						Val: "var2",
					},
					{
						IsString: true,
						Val: "string",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateDBGetReq(tt.args.argStrs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateDBGetReq() = %v, want %v", got, tt.want)
			}
		})
	}
}
