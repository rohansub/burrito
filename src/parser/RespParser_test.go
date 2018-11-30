package parser

import (
	"github.com/rcsubra2/burrito/src/db"
	"reflect"
	"testing"
)

func Test_createResp(t *testing.T) {
	type args struct {
		respStr string
	}
	tests := []struct {
		name    string
		args    args
		want    Resp
		wantErr bool
	}{
		{
			name: "File data test",
			args: args{
				respStr: "'hello.html'",
			},
			want: Resp{
				RespType: "FILE",
				Body:     "hello.html",
			},
			wantErr: false,
		},
		{
			name: "Invalid data (no quotes in file name)",
			args: args{
				respStr: "hello.html",
			},
			want:    Resp{},
			wantErr: true,
		},
		{
			name: "String data",
			args: args{
				respStr: "s'Hello World'",
			},
			want: Resp{
				RespType: "STR",
				Body:     "Hello World",
			},
			wantErr: false,
		},
		{
			name: "JSON data",
			args: args{
				respStr: `{ "hello" : "world", "zesty": { "breakfast" : "burrito" } } `,
			},
			want: Resp{
				RespType: "JSON",
				Body:     map[string]interface{}{
					"hello": "world",
					"zesty": map[string]interface{}{
						"breakfast": "burrito",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "JSON data, invalid format",
			args: args{
				respStr: `{ "hello" :world", "zesty": { "breakfast" : "burrito" } } `,
			},
			want: Resp{},
			wantErr: true,
		},
		{
			name: "JSON data Single Item",
			args: args{
				respStr: `{ "hello" : "world" } `,
			},
			want: Resp{
				RespType: "JSON",
				Body:     map[string]interface{}{
					"hello": "world",
				},
			},
			wantErr: false,
		},
		{
			name: "DB data",
			args: args{
				respStr: `DB.GET(varname,'hello',)`,
			},
			want: Resp{
				RespType: "DB",
				DBReq:    &db.GetReq{
					ArgNames: []db.Param{
						db.Param{
							false,
							"varname",
						},
						db.Param{
							true,
							"hello",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createResp(tt.args.respStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("createResp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createResp() = %v, want %v", got, tt.want)
			}
		})
	}
}
