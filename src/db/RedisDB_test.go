package db

import (
	"reflect"
	"testing"

	"github.com/rcsubra2/burrito/src/mockredis"
	"github.com/rcsubra2/burrito/src/utils"
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
		{
			name: "Test get multiple keys, all in db",
			fields: fields{
				db: mockredis.NewMockRedisClient(map[string]string{
					"jello":   "gorld",
					"goodbye": "world",
				}),
			},
			args: args{
				[]string{"jello", "goodbye"},
			},
			want: map[string]string{
				"jello":   "gorld",
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
			args: args{
				[]string{"jello", "goodbye"},
			},
			want: map[string]string{
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

func TestRedisDB_Delete(t *testing.T) {
	type fields struct {
		db RedisDBInterface
	}
	type args struct {
		keys    []string
		keysGet []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantGet map[string]string
	}{
		{
			name: "Delete multiple",
			fields: fields{
				db: mockredis.NewMockRedisClient(map[string]string{
					"hello": "world",
					"zesty": "burrito",
				}),
			},
			args: args{
				keys:    []string{"hello", "zesty"},
				keysGet: []string{"hello", "zesty"},
			},
			want:    true,
			wantGet: map[string]string{},
		},
		{
			name: "Delete single",
			fields: fields{
				db: mockredis.NewMockRedisClient(map[string]string{
					"hello": "world",
					"zesty": "burrito",
				}),
			},
			args: args{
				keys:    []string{"hello"},
				keysGet: []string{"hello", "zesty"},
			},
			want: true,
			wantGet: map[string]string{
				"zesty": "burrito",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &RedisDB{
				db: tt.fields.db,
			}
			if got := rc.Delete(tt.args.keys); got != tt.want {
				t.Errorf("RedisDB.Delete() = %v, want %v", got, tt.want)
			}

			if got2 := rc.Get(tt.args.keysGet); !reflect.DeepEqual(got2, tt.wantGet) {
				t.Errorf("Get after Del: RedisDB.Get() = %v, want %v", got2, tt.wantGet)
			}
		})
	}
}

func TestRedisDB_Set(t *testing.T) {
	type fields struct {
		db RedisDBInterface
	}
	type args struct {
		items []utils.Pair
		keysGet []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		wantGet map[string]string
	}{
		{
			name: "Set multiple",
			fields: fields{
				db: mockredis.NewMockRedisClient(map[string]string{}),
			},
			args: args{
				items:    []utils.Pair{
					{
						Fst: "hello",
						Snd: "world",
					},
					{
						Fst: "zesty",
						Snd: "burrito",
					},
				},
				keysGet: []string{"hello", "zesty"},
			},
			want: true,
			wantGet: map[string]string{
				"zesty": "burrito",
				"hello":"world",
			},
		},
		{
			name: "Set when value exists",
			fields: fields{
				db: mockredis.NewMockRedisClient(map[string]string{
					"hello": "slug",
				}),
			},
			args: args{
				items:    []utils.Pair{
					{
						Fst: "hello",
						Snd: "world",
					},
					{
						Fst: "zesty",
						Snd: "burrito",
					},
				},
				keysGet: []string{"hello", "zesty"},
			},
			want: true,
			wantGet: map[string]string{
				"zesty": "burrito",
				"hello":"world",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &RedisDB{
				db: tt.fields.db,
			}
			if got := rc.Set(tt.args.items); got != tt.want {
				t.Errorf("RedisDB.Set() = %v, want %v", got, tt.want)
			}
			if got2 := rc.Get(tt.args.keysGet); !reflect.DeepEqual(got2, tt.wantGet) {
				t.Errorf("RedisDB.Set() = %v, want %v", got2, tt.wantGet)
			}

		})
	}
}
