package config

import (
	"reflect"
	"testing"
)

func TestNewConfigFromFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "Test config file parsing",
			args: args{
				filename: "burritoconfig.json",
			},
			want: &Config{
				Name: "My Server",
				Databases: map[string]DbMeta{
					"rds": {
						DbType: "Redis",
						IsMock: true,
						Server: ServerMeta{
							Url: "localhost:9000",
							Password: "",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfigFromFile(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfigFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfigFromFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

