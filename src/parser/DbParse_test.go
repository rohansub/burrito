package parser

import (
	db2 "github.com/rcsubra2/burrito/src/redis"
	"reflect"
	"testing"

	"github.com/rcsubra2/burrito/src/db"
)


func Test_createDBCall(t *testing.T) {
	type args struct {
		respStr string
		databases    map[string]db.Database
	}
	tests := []struct {
		name string
		args args
		want *db.DatabaseAction
	}{
		{
			name: "Test database exists",
			args: args {
				respStr: "rds.GET(hello,)",
				databases: map[string]db.Database{
					"rds": db2.NewRedisDatabase(true, "", ""),
				},
			},
			want: &db.DatabaseAction{
				Name: "rds",
				Fname: "GET",
				Args: "hello,",
			},
		},
		{
			name: "Test database doesn't exist",
			args: args {
				respStr: "rds.GET(hello,)",
				databases: map[string]db.Database{
					"rds1": db2.NewRedisDatabase(true, "", ""),
				},
			},
			want: nil,
		},
		{
			name: "Test not correct db function syntax",
			args: args {
				respStr: "rdsGET(hello,)",
				databases: map[string]db.Database{
					"rds1": db2.NewRedisDatabase(true, "", ""),
				},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createDBCall(tt.args.respStr, tt.args.databases); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createDBCall() = %v, want %v", got, tt.want)
			}
		})
	}
}
