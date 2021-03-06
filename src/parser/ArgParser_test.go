package parser

import (
	"reflect"
	"testing"
)

func Test_createArg(t *testing.T) {
	type args struct {
		argStr string
	}
	tests := []struct {
		name    string
		args    args
		want    Arg
		wantErr bool
	}{
		{
			name: "Too many arguments",
			args: args{
				argStr: "('/zesty', 'GET', 'chocolate')",
			},
			want:    Arg{},
			wantErr: true,
		},
		{
			name: "Invalid path format",
			args: args{
				argStr: "('/zesty**', 'GET')",
			},
			want:    Arg{},
			wantErr: true,
		},
		{
			name: "Correct path, with GET explicitly written",
			args: args{
				argStr: "('/zesty/server', 'GET')",
			},
			want: Arg{
				path:    "/zesty/server",
				reqType: "GET",
			},
			wantErr: false,
		},
		{
			name: "Correct path, with PUT explicitly written",
			args: args{
				argStr: "('/zesty/server', 'PUT')",
			},
			want: Arg{
				path:    "/zesty/server",
				reqType: "PUT",
			},
			wantErr: false,
		},
		{
			name: "Correct path, with POST explicitly written",
			args: args{
				argStr: "('/zesty/server', 'POST')",
			},
			want: Arg{
				path:    "/zesty/server",
				reqType: "POST",
			},
			wantErr: false,
		},
		{
			name: "Correct path, with DELETE explicitly written",
			args: args{
				argStr: "('/zesty/server', 'DELETE')",
			},
			want: Arg{
				path:    "/zesty/server",
				reqType: "DELETE",
			},
			wantErr: false,
		},
		{
			name: "Correct path, invalid request",
			args: args{
				argStr: "('/zesty/server', 'XXX')",
			},
			want:    Arg{},
			wantErr: true,
		},
		{
			name: "Correct path, default request",
			args: args{
				argStr: "('/zest')",
			},
			want: Arg{
				path:    "/zest",
				reqType: "GET",
			},
			wantErr: false,
		},
		{
			name: "Default path, default request",
			args: args{
				argStr: "()",
			},
			want: Arg{
				path:    "/",
				reqType: "GET",
			},
			wantErr: false,
		},
		{
			name: "Path with parameter in path",
			args: args{
				argStr: "('/zesty/:server', 'GET')",
			},
			want: Arg{
				path:    "/zesty/:server",
				reqType: "GET",
			},
			wantErr: false,
		},
		{
			name: "Path with parameter in path, int type specified",
			args: args{
				argStr: "('/zesty/:server:int', 'GET')",
			},
			want: Arg{
				path:    "/zesty/:server:int",
				reqType: "GET",
			},
			wantErr: false,
		},
		{
			name: "Path with parameter in path, float type specified",
			args: args{
				argStr: "('/zesty/:server:flt', 'GET')",
			},
			want: Arg{
				path:    "/zesty/:server:flt",
				reqType: "GET",
			},
			wantErr: false,
		},
		{
			name: "Path with parameter in path, string type specified",
			args: args{
				argStr: "('/zesty/:server:str', 'GET')",
			},
			want: Arg{
				path:    "/zesty/:server:str",
				reqType: "GET",
			},
			wantErr: false,
		},
		{
			name: "Path with parameter in path, bad type specified",
			args: args{
				argStr: "('/zesty/:server:bogus', 'GET')",
			},
			want: Arg{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createArg(tt.args.argStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("createArg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createArg() = %v, want %v", got, tt.want)
			}
		})
	}
}
