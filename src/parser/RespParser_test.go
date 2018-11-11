package parser

import (
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
		// TODO: Add test cases.
		{
			name: "File data test",
			args: args{
				respStr: "'hello.html'",
			},
			want: Resp{
				respType: "FILE",
				body:     "hello.html",
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
				respType: "STR",
				body:     "Hello World",
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
