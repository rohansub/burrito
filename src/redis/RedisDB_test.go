package redis

import (
	"testing"
)

//func TestCreateGet(t *testing.T) {
//	url := environment.CreateEnv()
//	url.Add(*environment.CreateStringEntry("hello", "stuff"))
//	url.Add(*environment.CreateStringEntry("channa", "sambar"))
//
//	form := environment.CreateEnv()
//	resp := environment.CreateEnv()
//
//	type args struct {
//		Client RedisDBClientInterface
//		args   string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		group   environment.EnvironmentGroup
//		want    db.DatabaseFunction
//		wantErr bool
//	}{
//		{
//			name: "Test create Get, variables and strings",
//			args: args{
//				args:   "hello, 'world', channa, 'masala',",
//				Client: redis.NewMockRedisClient(map[string]string{}),
//			},
//			group: *environment.CreateEnvironmentGroup(url, form, resp),
//			want: func(group environment.EnvironmentGroup) (map[string]interface{}, error) {
//				return Get([]string{"hello", "'world'", "channa", "'masala'"},
//					redis.NewMockRedisClient(map[string]string{}),
//					group), nil
//			},
//			wantErr: false,
//		},
//		{
//			name: "Test create Get, variables and strings, with db values",
//			args: args{
//				args: "hello, 'world', channa, 'masala',",
//				Client: redis.NewMockRedisClient(map[string]string{
//					"stuff":  "people",
//					"sambar": "dosa",
//					"world":  "things",
//					"masala": "chutney",
//				}),
//			},
//			group: *environment.CreateEnvironmentGroup(url, form, resp),
//			want: func(group environment.EnvironmentGroup) (map[string]interface{}, error) {
//				return Get([]string{"hello", "'world'", "channa", "'masala'"},
//					redis.NewMockRedisClient(map[string]string{
//						"stuff":  "people",
//						"sambar": "dosa",
//						"world":  "things",
//						"masala": "chutney",
//					}),
//					group), nil
//			},
//
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := NewGetFunction(tt.args.Client, tt.args.args)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("NewGetFunction() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//
//			got2, _ := got(tt.group)
//			want2, _ := tt.want(tt.group)
//			if !reflect.DeepEqual(got2, want2) {
//				t.Errorf("NewGetFunction(tt.group) = %v, want %v", got2, want2)
//			}
//
//		})
//	}
//}
//
//func TestGet(t *testing.T) {
//	url := environment.CreateEnv()
//	url.Add(*environment.CreateStringEntry("zesty", "burrito"))
//	url.Add(*environment.CreateStringEntry("snack", "snake"))
//
//	form := environment.CreateEnv()
//	resp := environment.CreateEnv()
//
//	type args struct {
//		keys  []string
//		db    RedisDBClientInterface
//		group environment.EnvironmentGroup
//	}
//	tests := []struct {
//		name string
//		args args
//		want map[string]interface{}
//	}{
//		{
//			name: "Test get strings and variables",
//			args: args{
//				keys: []string{
//					"'chocolate'",
//					"snack",
//					"zesty",
//				},
//				db: redis.NewMockRedisClient(map[string]string{
//					"chocolate": "cake",
//					"snake":     "evil",
//					"burrito":   "supreme",
//				}),
//				group: *environment.CreateEnvironmentGroup(url, form, resp),
//			},
//			want: map[string]interface{}{
//				"chocolate": "cake",
//				"snake":     "evil",
//				"burrito":   "supreme",
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := Get(tt.args.keys, tt.args.db, tt.args.group); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Get() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestNewSetFunction(t *testing.T) {
//	type args struct {
//		Client RedisDBClientInterface
//		args   string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    db.DatabaseFunction
//		wantErr bool
//	}{
//		{
//			name: "Test create Set, variables and strings",
//			args: args{
//				args:   "(hello, 'world'), (channa, 'masala'),",
//				Client: redis.NewMockRedisClient(map[string]string{}),
//			},
//			want:    nil,
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := NewSetFunction(tt.args.Client, tt.args.args)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("NewSetFunction() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("NewSetFunction() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}




func TestRedisDatabase_IsCorrectSyntax(t *testing.T) {
	type fields struct {
		client RedisDBClientInterface
	}
	type args struct {
		fname string
		args  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Test unrecognized function",
			fields: fields{
				client: NewRedisDatabase(true, "", "").Client,
			},
			args: args {
				fname: "Fakefunc",
				args: "zesty, 'Hello',",
			},
			want: false,
		},
		{
			name: "Test GET syntax correct",
			fields: fields{
				client: NewRedisDatabase(true, "", "").Client,
			},
			args: args {
				fname: "GET",
				args: "zesty, 'Hello',",
			},
			want: true,
		},
		{
			name: "Test GET syntax incorrect",
			fields: fields{
				client: NewRedisDatabase(true, "", "").Client,
			},
			args: args {
				fname: "GET",
				args: "zesty, 'Hello'",
			},
			want: false,
		},
		{
			name: "Test DELETE syntax correct",
			fields: fields{
				client: NewRedisDatabase(true, "", "").Client,
			},
			args: args {
				fname: "DEL",
				args: "zesty, 'Hello',",
			},
			want: true,
		},
		{
			name: "Test DEL syntax incorrect",
			fields: fields{
				client: NewRedisDatabase(true, "", "").Client,
			},
			args: args {
				fname: "DEL",
				args: "zesty, 'Hello'",
			},
			want: false,
		},
		{
			name: "Test SET syntax correct",
			fields: fields{
				client: NewRedisDatabase(true, "", "").Client,
			},
			args: args {
				fname: "SET",
				args: "(zesty, channa), ('masala', burrito),",
			},
			want: true,
		},
		{
			name: "Test SET syntax incorrect",
			fields: fields{
				client: NewRedisDatabase(true, "", "").Client,
			},
			args: args {
				fname: "SET",
				args: "zesty, 'Hello',",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rd := &RedisDatabase{
				Client: tt.fields.client,
			}
			if got := rd.IsCorrectSyntax(tt.args.fname, tt.args.args); got != tt.want {
				t.Errorf("RedisDatabase.IsCorrectSyntax() = %v, want %v", got, tt.want)
			}
		})
	}
}
