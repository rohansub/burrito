package db

import (
	"github.com/rcsubra2/burrito/src/mockredis"
	"reflect"
	"testing"

	"github.com/rcsubra2/burrito/src/environment"
)

func TestNewRedisDB(t *testing.T) {
	type args struct {
		client RedisDBInterface
	}
	tests := []struct {
		name string
		args args
		want *RedisDB
	}{
		{
			name: "Create RedisDB",
			args: args{
				client:mockredis.NewMockRedisClient(map[string]string{
					"key": "value",
				}),
			},
			want: &RedisDB{
				db: mockredis.NewMockRedisClient(map[string]string{
					"key": "value",
				}),
			},

		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRedisDB(tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRedisDB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisDB_Get(t *testing.T) {
	env := environment.CreateEnv()
	env.Add(*environment.CreateStringEntry("hello", "goodbye"))

	type fields struct {
		db RedisDBInterface
	}
	type args struct {
		req  GetReq
		envs []*environment.Env
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]string
	}{
		{
			name: "GET, field exists",
			fields: fields{
				db: mockredis.NewMockRedisClient(map[string]string{
					"jello": "gorld",
					"hello": "world",
				}),
			},
			args: args {
				req: GetReq{
					ArgNames: []Param {
						{
							IsString: true,
							Val: "jello",
						},
						{
							IsString: true,
							Val: "hello",
						},
					},
				},
				envs: []*environment.Env{},
			},
			want: map[string]string {
				"jello": "gorld",
				"hello": "world",
			},
		},
		{
			name: "GET, field does not exist",
			fields: fields{
				db: mockredis.NewMockRedisClient(map[string]string{
					"jello": "gorld",
				}),
			},
			args: args {
				req: GetReq{
					ArgNames: []Param {
						{
							IsString: true,
							Val: "jello",
						},
						{
							IsString: true,
							Val: "hello",
						},
					},
				},
				envs: []*environment.Env{},
			},
			want: map[string]string {
				"jello": "gorld",
			},
		},
		{
			name: "GET, use environment field",
			fields: fields{
				db: mockredis.NewMockRedisClient(map[string]string{
					"jello": "gorld",
					"goodbye": "world",
				}),
			},
			args: args {
				req: GetReq{
					ArgNames: []Param {
						{
							IsString: true,
							Val: "jello",
						},
						{
							IsString: false,
							Val: "hello",
						},
					},
				},
				envs: []*environment.Env{env},
			},
			want: map[string]string {
				"jello": "gorld",
				"goodbye": "world",
			},
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &RedisDB{
				db: tt.fields.db,
			}
			if got := rc.Get(tt.args.req, tt.args.envs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedisDB.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
