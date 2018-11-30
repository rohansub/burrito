package db

import (
	"reflect"
	"testing"

	"github.com/rcsubra2/burrito/src/mockredis"
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
				client: mockredis.NewMockRedisClient(map[string]string{
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
	type fields struct {
		db RedisDBInterface
	}
	type args struct {
		keys []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]string
	}{
		// TODO: Add test cases.
		{
			name: "Test get multiple keys, all in db",
			fields: fields{
				db: mockredis.NewMockRedisClient(map[string]string{
					"jello": "gorld",
					"goodbye": "world",
				}),
			},
			args: args {
				[]string{"jello", "goodbye"},
			},
			want: map[string]string {
				"jello": "gorld",
				"goodbye": "world",
			},
		},
		{
			name: "Test get multiple keys, not all in db",
			fields: fields{
				db: mockredis.NewMockRedisClient(map[string]string{
					"jello": "gorld",
					"hello": "world",
				}),
			},
			args: args {
				[]string{"jello", "goodbye"},
			},
			want: map[string]string {
				"jello": "gorld",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &RedisDB{
				db: tt.fields.db,
			}
			if got := rc.Get(tt.args.keys); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedisDB.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}


