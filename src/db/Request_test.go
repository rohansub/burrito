package db

import (
	"reflect"
	"testing"

	"github.com/rcsubra2/burrito/src/environment"
	"github.com/rcsubra2/burrito/src/mockredis"
	"github.com/rcsubra2/burrito/src/utils"
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
		wantOk bool
	}{
		{
			name: "Test String",
			fields: fields{
				IsString: true,
				Val:      "hello",
			},
			args: args{
				envs: []*environment.Env{env1, env2},
			},
			want:   "hello",
			wantOk: true,
		},
		{
			name: "Test Variable",
			fields: fields{
				IsString: false,
				Val:      "hello",
			},
			args: args{
				envs: []*environment.Env{env1, env2},
			},
			want:   "world",
			wantOk: true,
		},
		{
			name: "Test Variable Not in entry",
			fields: fields{
				IsString: false,
				Val:      "nothere",
			},
			args: args{
				envs: []*environment.Env{env1, env2},
			},
			want:   "",
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
			if (ok != false) != tt.wantOk {
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
				argStrs: []string{"'hello'", "variable", "var2", "'string'"},
			},
			want: &GetReq{
				ArgNames: []Param{
					{
						IsString: true,
						Val:      "hello",
					},
					{
						IsString: false,
						Val:      "variable",
					},
					{
						IsString: false,
						Val:      "var2",
					},
					{
						IsString: true,
						Val:      "string",
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

func TestGetReq_Run(t *testing.T) {
	env := environment.CreateEnv()
	env.Add(*environment.CreateStringEntry("Maize", "Burrito"))

	type fields struct {
		ArgNames []Param
	}
	type args struct {
		client Database
		envs   []*environment.Env
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]string
		err    error
	}{
		{
			name: "Test Get request all strings",
			fields: fields{
				ArgNames: []Param{
					{
						IsString: true,
						Val:      "Veggie",
					},
					{
						IsString: true,
						Val:      "Burrito",
					},
				},
			},
			args: args{
				client: NewRedisDB(
					mockredis.NewMockRedisClient(map[string]string{
						"Veggie":  "great",
						"Burrito": "food",
					}),
				),
				envs: []*environment.Env{},
			},
			want: map[string]string{
				"Veggie":  "great",
				"Burrito": "food",
			},
		},
		{
			name: "Test Get request strings and env variables",
			fields: fields{
				ArgNames: []Param{
					{
						IsString: true,
						Val:      "Veggie",
					},
					{
						IsString: false,
						Val:      "Maize",
					},
				},
			},
			args: args{
				client: NewRedisDB(
					mockredis.NewMockRedisClient(map[string]string{
						"Veggie":  "great",
						"Burrito": "food",
					}),
				),
				envs: []*environment.Env{
					env,
				},
			},
			want: map[string]string{
				"Veggie":  "great",
				"Burrito": "food",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &GetReq{
				ArgNames: tt.fields.ArgNames,
			}
			got, err := r.Run(tt.args.client, tt.args.envs)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Err = %v, want %v", err, tt.err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetReq.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateDBSetReq(t *testing.T) {
	type args struct {
		argStrs []string
	}
	tests := []struct {
		name string
		args args
		want *SetReq
	}{
		{
			name: "Test Creating SET request variable and string",
			args: args{
				argStrs: []string{"('hello', variable)", "(var2,'string')"},
			},
			want: &SetReq{
				ArgNames: []utils.Pair{
					{
						Fst: Param{
							IsString: true,
							Val:      "hello",
						},
						Snd: Param{
							IsString: false,
							Val:      "variable",
						},
					},
					{
						Fst: Param{
							IsString: false,
							Val:      "var2",
						},
						Snd: Param{
							IsString: true,
							Val:      "string",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateDBSetReq(tt.args.argStrs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateDBSetReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetReq_Run(t *testing.T) {
	env := environment.CreateEnv()
	env.Add(*environment.CreateStringEntry("variable", "world"))
	env.Add(*environment.CreateStringEntry("var2", "ThisIs"))


	type fields struct {
		ArgNames []utils.Pair
	}
	type args struct {
		client Database
		envs   []*environment.Env
	}
	type getFields struct {
		ArgNames []Param
	}
	tests := []struct {
		name    	string
		fields  	fields
		args    	args
		want    	map[string]string
		wantErr 	bool
		getFields 	getFields
		wantGet  	map[string] string
		wantGetErr 	bool
	}{
		{
			name: "Test Set request Run - strings only",
			fields: fields {
				ArgNames: []utils.Pair{
					{
						Fst: Param{
							IsString: true,
							Val:      "hello",
						},
						Snd: Param{
							IsString: true,
							Val:      "variable",
						},
					},
					{
						Fst: Param{
							IsString: true,
							Val:      "var2",
						},
						Snd: Param{
							IsString: true,
							Val:      "string",
						},
					},
				},
			},
			args: args {
				client: NewRedisDB(
					mockredis.NewMockRedisClient(map[string]string{}),
				),
				envs: []*environment.Env{},
			},
			want: map[string]string{},
			wantErr: false,
			getFields: getFields{
				ArgNames: []Param{
					{
						IsString: true,
						Val:      "hello",
					},
					{
						IsString: true,
						Val:      "var2",
					},
				},
			},
			wantGet: map[string]string {
				"hello": "variable",
				"var2": "string",
			},
			wantGetErr: false,
		},
		{
			name: "Test Set request Run - strings and variables",
			fields: fields {
				ArgNames: []utils.Pair{
					{
						Fst: Param{
							IsString: true,
							Val:      "hello",
						},
						Snd: Param{
							IsString: false,
							Val:      "variable",
						},
					},
					{
						Fst: Param{
							IsString: false,
							Val:      "var2",
						},
						Snd: Param{
							IsString: true,
							Val:      "string",
						},
					},
				},
			},
			args: args {
				client: NewRedisDB(
					mockredis.NewMockRedisClient(map[string]string{}),
				),
				envs: []*environment.Env{env},
			},
			want: map[string]string{},
			wantErr: false,
			getFields: getFields{
				ArgNames: []Param{
					{
						IsString: true,
						Val:      "hello",
					},
					{
						IsString: false,
						Val:      "var2",
					},
				},
			},
			wantGet: map[string]string {
				"hello": "world",
				"ThisIs": "string",
			},
			wantGetErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &SetReq{
				ArgNames: tt.fields.ArgNames,
			}
			got, err := req.Run(tt.args.client, tt.args.envs)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetReq.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetReq.Run() = %v, want %v", got, tt.want)
			}

			gr := &GetReq{
				ArgNames: tt.getFields.ArgNames,
			}

			got2, err := gr.Run(tt.args.client, tt.args.envs)
			if (err != nil) != tt.wantGetErr {
				t.Errorf("Get after Set failed: Err = %v, want %v", err, tt.wantGetErr)
			}
			if !reflect.DeepEqual(got2, tt.wantGet) {
				t.Errorf("Get after Set failed: GetReq.Run() = %v, want %v", got2, tt.wantGet)
			}
		})
	}
}
