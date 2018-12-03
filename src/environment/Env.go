package environment

// Env - Consists of key value pair, where key is a variable name, and value is an EnvEntry corresponding
// to that variable
type Env struct{
	data map[string]EnvEntry
}

// CreateEnv - Create and return an object of type Env
func CreateEnv() *Env{
	d := make(map[string]EnvEntry)
	return &Env {
		data: d,
	}
}

// Add - add an entry to the environment, only works for non-nil entry
func (e * Env) Add(entry EnvEntry) {
	e.data[entry.name] = entry
}

func (e * Env) Get(entryName string) *EnvEntry{
	entry, ok := e.data[entryName]
	if !ok {
		return nil
	}
	return &entry
}

func (e * Env) GetValue(entryName string) interface{}{
	entry, ok := e.data[entryName]
	if !ok {
		return nil
	}
	if entry.isStr {
		return entry.valStr
	} else if entry.isFlt {
		return entry.valFlt
	} else { // otherwise it is an int
		return entry.valInt
	}
}


func (e * Env) Data() map[string]interface{}{
	d := make(map[string]interface{})
	for k, _ := range e.data {
		d[k] = e.GetValue(k)
	}
	return d
}


// EnvEntry - Consists of an Variable name, as well as a value. Value can be a string,
// float64, or int64
type EnvEntry struct {
	name   string
	valInt int64
	valFlt float64
	valStr string
	isInt bool
	isFlt bool
	isStr bool
}

// CreateIntEntry - create an integer entry for an Environment.
func CreateIntEntry(name string, i int64) *EnvEntry {
	return &EnvEntry{
		name: name,
		valInt: i,
		isInt: true,
	}
}

// CreateFloatEntry - create an integer entry for an Environment.
func CreateFloatEntry(name string, f float64) *EnvEntry {
	return &EnvEntry{
		name: name,
		valFlt: f,
		isFlt: true,
	}
}

// CreateStringEntry - create an integer entry for an Environment.
func CreateStringEntry(name string, st string) *EnvEntry {
	return &EnvEntry{
		name: name,
		valStr: st,
		isStr: true,
	}
}

// BurritoTemplateData - Struct to organize burrito data into Template-parasable data
type BurritoTemplateData struct {
	Url interface{}
	Form interface{}
	Data interface{}
}


func CreateBurritoTemplateData(urlEnv * Env, respEnv * Env) BurritoTemplateData{
	t := BurritoTemplateData{}

	t.Url = urlEnv.Data()
	t.Data = respEnv.Data()
	return t
}


// EnvironmentGroup - Group of environment structures used when a request is parsed
type EnvironmentGroup struct {
	Url	*Env
	Form *Env
	Resp *Env
}

// CreateEnvironmentGroup - Creates and environment group
func CreateEnvironmentGroup(url *Env, form *Env, resp *Env) *EnvironmentGroup{
	return &EnvironmentGroup{
		Url: url,
		Form: form,
		Resp: resp,
	}
}


// GetValue - get value specified in an environment group
func (eg * EnvironmentGroup) GetValue(key string) interface{}{
	inResp := eg.Url.GetValue(key)
	if inResp != nil {
		return inResp
	}

	inForm := eg.Url.GetValue(key)
	if inForm != nil {
		return inForm
	}

	inUrl := eg.Url.GetValue(key)
	if inUrl != nil {
		return inUrl
	}

	return nil
}

func (eg * EnvironmentGroup) Dump() interface{} {
	return BurritoTemplateData{
		Url: eg.Url.Data(),
		Form: eg.Form.Data(),
		Data: eg.Resp.Data(),
	}
}