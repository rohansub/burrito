package environment

import (
	"reflect"
	"testing"
)

func TestCreateEnv(t *testing.T) {
	tests := []struct {
		name string
		want *Env
	}{
		{
			name: "Test Creating Empty Env",
			want: &Env{
				data: make(map[string]EnvEntry),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateEnv(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_Add(t *testing.T) {
	type fields struct {
		data map[string]EnvEntry
	}
	type args struct {
		entry EnvEntry
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		fieldsAfter fields
	}{
		{
			name: "Add field when variable not in Env",
			fields: fields{
				data: map[string]EnvEntry{},
			},
			args: args{
				entry: *CreateIntEntry("name", 5),
			},
			fieldsAfter: fields{
				data: map[string]EnvEntry{
					"name": *CreateIntEntry("name", 5),
				},
			},
		},
		{
			name: "Add field when variable is already in Env",
			fields: fields{
				data: map[string]EnvEntry{
					"name": *CreateStringEntry("name", "old"),
				},
			},
			args: args{
				entry: *CreateIntEntry("name", 5),
			},
			fieldsAfter: fields{
				data: map[string]EnvEntry{
					"name": *CreateIntEntry("name", 5),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Env{
				data: tt.fields.data,
			}
			e.Add(tt.args.entry)
			got := e.data
			want := tt.fieldsAfter.data
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Env().data = %v, want %v", got, want)
			}
		})
	}
}

func TestEnv_Get(t *testing.T) {
	type fields struct {
		data map[string]EnvEntry
	}
	type args struct {
		entryName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *EnvEntry
	}{
		{
			name: "Test Get, item in environment",
			fields: fields{
				data: map[string]EnvEntry{
					"zesty": *CreateStringEntry("zesty", "burrito"),
				},
			},
			args: args{
				"zesty",
			},
			want: CreateStringEntry("zesty", "burrito"),
		},
		{
			name: "Test Get, item not found in environment",
			fields: fields{
				data: map[string]EnvEntry{},
			},
			args: args{
				"zesty",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Env{
				data: tt.fields.data,
			}
			if got := e.Get(tt.args.entryName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Env.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateIntEntry(t *testing.T) {
	type args struct {
		name string
		i    int64
	}
	tests := []struct {
		name string
		args args
		want *EnvEntry
	}{
		{
			name: "Test Create Int Entry",
			args: args{
				name: "siesta",
				i:    0xBAADF00D,
			},
			want: &EnvEntry{
				name:   "siesta",
				isInt:  true,
				valInt: 0xBAADF00D,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateIntEntry(tt.args.name, tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateIntEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateFloatEntry(t *testing.T) {
	type args struct {
		name string
		f    float64
	}
	tests := []struct {
		name string
		args args
		want *EnvEntry
	}{
		{
			name: "Test Create Float Entry",
			args: args{
				name: "siesta",
				f:    66.66,
			},
			want: &EnvEntry{
				name:   "siesta",
				isFlt:  true,
				valFlt: 66.66,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateFloatEntry(tt.args.name, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateFloatEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateStringEntry(t *testing.T) {
	type args struct {
		name string
		st   string
	}
	tests := []struct {
		name string
		args args
		want *EnvEntry
	}{
		{
			name: "Test Create Float Entry",
			args: args{
				name: "siesta",
				st:   "fiesta",
			},
			want: &EnvEntry{
				name:   "siesta",
				isStr:  true,
				valStr: "fiesta",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateStringEntry(tt.args.name, tt.args.st); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateStringEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateBurritoTemplateData(t *testing.T) {
	type args struct {
		urlEnv  *Env
		respEnv *Env
	}
	tests := []struct {
		name string
		args args
		want BurritoTemplateData
	}{
		{
			name: "Create Empty Burrito Template",
			args: args{
				urlEnv: &Env{
					map[string]EnvEntry{},
				},
				respEnv: &Env{
					map[string]EnvEntry{},
				},
			},
			want: BurritoTemplateData{
				Url: (&Env{
					map[string]EnvEntry{},
				}).Data(),
				Data: (&Env{
					map[string]EnvEntry{},
				}).Data(),
			},
		},
		{
			name: "Create Burrito Template - contains data",
			args: args{
				urlEnv: &Env{
					map[string]EnvEntry{
						"hello": *CreateIntEntry("hello", 1),
						"zesty": *CreateStringEntry("zesty", "burrito"),
						"siesta": *CreateFloatEntry("siesta", 999.9),
					},
				},
				respEnv: &Env{
					map[string]EnvEntry{
						"rhello": *CreateIntEntry("rhello", 1),
						"rzesty": *CreateStringEntry("rzesty", "burrito"),
						"rsiesta": *CreateFloatEntry("rsiesta", 999.9),
					},
				},
			},
			want: BurritoTemplateData{
				Url: (&Env{
					map[string]EnvEntry{
						"hello": *CreateIntEntry("hello", 1),
						"zesty": *CreateStringEntry("zesty", "burrito"),
						"siesta": *CreateFloatEntry("siesta", 999.9),
					},
				}).Data(),
				Data: (&Env{
					map[string]EnvEntry{
						"rhello": *CreateIntEntry("rhello", 1),
						"rzesty": *CreateStringEntry("rzesty", "burrito"),
						"rsiesta": *CreateFloatEntry("rsiesta", 999.9),
					},
				}).Data(),
			},
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateBurritoTemplateData(tt.args.urlEnv, tt.args.respEnv); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateBurritoTemplateData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_GetValue(t *testing.T) {
	type fields struct {
		data map[string]EnvEntry
	}
	type args struct {
		entryName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		{
			name: "Get Value int",
			fields: fields {
				data: map[string]EnvEntry{
					"key": *CreateIntEntry("key", 4),
				},
			},
			args: args{"key"},
			want: int64(4),
		},
		{
			name: "Get Value str",
			fields: fields {
				data: map[string]EnvEntry{
					"key": *CreateStringEntry("key", "hello"),
				},
			},
			args: args{"key"},
			want: "hello",
		},
		{
			name: "Get Value float",
			fields: fields {
				data: map[string]EnvEntry{
					"key": *CreateFloatEntry("key", 4.2),
				},
			},
			args: args{"key"},
			want: 4.2,
		},
		{
			name: "Key not in Env",
			fields: fields {
				data: map[string]EnvEntry{
					"key": *CreateFloatEntry("key", 4.2),
				},
			},
			args: args{"not"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Env{
				data: tt.fields.data,
			}
			if got := e.GetValue(tt.args.entryName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Env.GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_Data(t *testing.T) {
	type fields struct {
		data map[string]EnvEntry
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]interface{}
	}{
		{
			name: "Extract data values",
			fields: fields {
				data: map[string]EnvEntry{
					"i": *CreateIntEntry("i", 1),
					"f": *CreateFloatEntry("f", 1.2),
					"st": *CreateStringEntry("st", "GG"),
				},
			},
			want: map[string]interface{}{
				"i": int64(1),
				"f": 1.2,
				"st": "GG",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Env{
				data: tt.fields.data,
			}
			if got := e.Data(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Env.Dump() = %v, want %v", got, tt.want)
			}
		})
	}
}
