package parser

import (
	"reflect"
	"testing"

	"github.com/rcsubra2/burrito/src/db"
)

func Test_createRespForDB(t *testing.T) {
	type args struct {
		respStr string
	}
	tests := []struct {
		name string
		args args
		want db.Req
	}{
		{
			name: "Test Strings",
			args: args {
				"DB.GET('hello','world',)",
			},
			want: &db.GetReq {
				ArgNames: []db.Param{
					{
						IsString: true,
						Val:      "hello",
					},
					{
						IsString: true,
						Val:      "world",
					},
				},
			},

		},
		{
			name: "Test Variables",
			args: args {
				"DB.GET(zesty, burrito,)",
			},
			want: &db.GetReq{
				ArgNames: []db.Param{
					{
						IsString: false,
						Val:      "zesty",
					},
					{
						IsString: false,
						Val:      "burrito",
					},
				},
			},

		},
		{
			name: "Test Variables an strings",
			args: args {
				"DB.GET(zesty, 'burrito', tomorrow,)",
			},
			want: &db.GetReq{
				ArgNames: []db.Param{
					{
						IsString: false,
						Val: "zesty",
					},
					{
						IsString: true,
						Val: "burrito",
					},
					{
						IsString: false,
						Val: "tomorrow",
					},

				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createRespForDB(tt.args.respStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createRespForDB() = %v, want %v", got, tt.want)
			}
		})
	}
}
