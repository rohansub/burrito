package redis

import (
	"reflect"
	"testing"
	"errors"
	"github.com/go-redis/redis"
	"time"
)

func TestNewMockRedisClient(t *testing.T) {
	type args struct {
		init map[string]string
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		{
			name: "Init Test",
			args: args {
				map[string]string {
					"key1": "value1",
					"key2": "value2",
				},
			},
			want: &Client{
				data: map[string]string {
					"key1": "value1",
					"key2": "value2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMockRedisClient(tt.args.init); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMockRedisClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Get(t *testing.T) {
	type fields struct {
		data map[string]string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *redis.StringCmd
	}{
		{
			name: "Test Get, key exists",
			fields: fields{
				map[string]string {
					"key": "val",
				},
			},
			args: args{
				"key",
			},
			want: redis.NewStringResult("val", nil),
		},
		{
			name: "Test Get, key doesn't exist",
			fields: fields{
				map[string]string {
					"key2": "val",
				},
			},
			args: args{
				"key",
			},
			want: redis.NewStringResult("", errors.New("Key not found")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				data: tt.fields.data,
			}
			if got := c.Get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Set(t *testing.T) {
	type fields struct {
		data map[string]string
	}
	type args struct {
		key   string
		value interface{}
		t	  time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantFields fields
	}{
		{
			name: "Test SET add new key",
			fields: fields {
				map[string]string{},
			},
			args: args {
				key: "key",
				value: "val",
				t: 0,
			},
			wantFields: fields {
				map[string]string{
					"key": "val",
				},
			},
		},
		{
			name: "Test SET modify value of old key",
			fields: fields {
				map[string]string{
					"key": "val1",
				},
			},
			args: args {
				key: "key",
				value: "val2",
				t: 0,
			},
			wantFields: fields {
				map[string]string{
					"key": "val2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				data: tt.fields.data,
			}
			c.Set(tt.args.key, tt.args.value, tt.args.t)
			if !reflect.DeepEqual(c.data, tt.wantFields.data) {
				t.Errorf("c.data = %v, want %v", c.data, tt.wantFields.data)
			}
		})
	}
}
