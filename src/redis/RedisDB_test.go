package db

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/rcsubra2/burrito/src/db"
	"github.com/rcsubra2/burrito/src/environment"
	"github.com/rcsubra2/burrito/src/mockredis"
)

func TestCreateGet(t *testing.T) {
	url := environment.CreateEnv()
	url.Add(*environment.CreateStringEntry("hello", "stuff"))
	url.Add(*environment.CreateStringEntry("channa", "sambar"))

	form := environment.CreateEnv()
	resp := environment.CreateEnv()



	type args struct {
		client RedisDBClientInterface
		args   string
	}
	tests := []struct {
		name    string
		args    args
		group   environment.EnvironmentGroup
		want    db.DatabaseFunction
		wantErr bool
	}{
		{
			name: "Test create Get, variables and strings",
			args: args{
				args:   "hello, 'world', channa, 'masala',",
				client: mockredis.NewMockRedisClient(map[string]string{}),
			},
			group: *environment.CreateEnvironmentGroup(url, form, resp),
			want: func(group environment.EnvironmentGroup) (interface{}, error) {
				return Get([]string{"hello", "'world'", "channa", "'masala'"},
					mockredis.NewMockRedisClient(map[string]string{}),
					group), nil
			},
			wantErr: false,
		},
		{
			name: "Test create Get, variables and strings, with db values",
			args: args{
				args: "hello, 'world', channa, 'masala',",
				client: mockredis.NewMockRedisClient(map[string]string{
					"stuff":  "people",
					"sambar": "dosa",
					"world":  "things",
					"masala": "chutney",
				}),
			},
			group: *environment.CreateEnvironmentGroup(url, form, resp),
			want: func(group environment.EnvironmentGroup) (interface{}, error) {
				return Get([]string{"hello", "'world'", "channa", "'masala'"},
					mockredis.NewMockRedisClient(map[string]string{
						"stuff":  "people",
						"sambar": "dosa",
						"world":  "things",
						"masala": "chutney",
					}),
					group), nil
			},

			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateGet(tt.args.client, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got2, _ := got(tt.group)
			want2, _ := tt.want(tt.group)
			fmt.Println(got2)
			if !reflect.DeepEqual(got2, want2) {
				t.Errorf("CreateGet(tt.group) = %v, want %v", got2, want2)
			}

		})
	}
}

func TestGet(t *testing.T) {
	url := environment.CreateEnv()
	url.Add(*environment.CreateStringEntry("zesty", "burrito"))
	url.Add(*environment.CreateStringEntry("snack", "snake"))

	form := environment.CreateEnv()
	resp := environment.CreateEnv()

	type args struct {
		keys  []string
		db    RedisDBClientInterface
		group environment.EnvironmentGroup
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "Test get strings and variables",
			args: args {
				keys: []string {
					"'chocolate'",
					"snack",
					"zesty",
				},
				db: mockredis.NewMockRedisClient(map[string]string{
					"chocolate":  "cake",
					"snake": "evil",
					"burrito":  "supreme",
				}),
				group: *environment.CreateEnvironmentGroup(url, form, resp),
			},
			want: map[string]string {
				"chocolate":  "cake",
				"snake": "evil",
				"burrito":  "supreme",
			},
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Get(tt.args.keys, tt.args.db, tt.args.group); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
