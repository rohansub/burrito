package mux


type Env struct{
	data map[string]EnvEntry
}

func CreateEnv() *Env{
	d := make(map[string]EnvEntry)
	return &Env {
		data: d,
	}
}


type EnvEntry struct {
	name   string
	valInt int64
	valFlt float64
	valStr string
	isInt bool
	isFlt bool
	isStr bool
}

func CreateIntEntry(name string, i int64) *EnvEntry {
	return &EnvEntry{
		name: name,
		valInt: i,
		isInt: true,
	}
}


func CreateFloatEntry(name string, f float64) *EnvEntry {
	return &EnvEntry{
		name: name,
		valFlt: f,
		isFlt: true,
	}
}


func CreateStringEntry(name string, st string) *EnvEntry {
	return &EnvEntry{
		name: name,
		valStr: st,
		isStr: true,
	}
}
