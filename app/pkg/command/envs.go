package command

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"sync"
)

type Envs struct {
	kvs map[string]string
	mux *sync.RWMutex
}

func NewEnvs() *Envs {
	return &Envs{
		kvs: make(map[string]string),
		mux: &sync.RWMutex{},
	}
}

func NewEnvsByMapStringAny(source map[string]any) *Envs {
	res := NewEnvs()
	for k, v := range source {
		res.kvs[k] = fmt.Sprintf("%v", v)
	}
	return res
}

func NewEnvsBySliceKV(source []string) *Envs {
	res := NewEnvs()
	for _, v := range source {
		r := strings.SplitN(v, "=", 2)
		if len(r) == 2 {
			res.kvs[strings.Trim(r[0], " ")] = strings.Trim(r[1], " ")
		}
	}
	return res
}

func (e *Envs) SliceKV() []string {
	e.mux.RLock()
	defer e.mux.RUnlock()
	res := make([]string, len(e.kvs))
	idx := 0
	for k, v := range e.kvs {
		res[idx] = fmt.Sprintf("%s=\"%s\"", k, v)
		idx++
	}
	return res
}

func (e *Envs) Add(k string, v any) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.kvs[k] = fmt.Sprintf("%v", v)
}

// Pick 选择一些key组成新的envs
func (e *Envs) Pick(keys ...string) *Envs {
	e.mux.RLock()
	defer e.mux.RUnlock()
	res := NewEnvs()
	for _, k := range keys {
		if v, ok := e.kvs[k]; ok {
			res.kvs[k] = v
		}
	}
	return res
}

func (e *Envs) MapString() map[string]string {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.kvs
}

func (e *Envs) String() string {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return strings.Join(e.SliceKV(), " ")
}

func (e *Envs) Empty() bool {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return len(e.kvs) == 0
}

func (e Envs) Value() (driver.Value, error) {
	return strings.Join(e.SliceKV(), ";"), nil
}

func (e *Envs) Scan(value interface{}) error {
	var str string
	switch value := value.(type) {
	case []byte:
		str = string(value)
		break
	case string:
		str = value
		break
	default:
		return fmt.Errorf("unable to scan %T into Envs", value)
	}
	if str == "" {
		e = NewEnvs()
	} else {
		e = NewEnvsBySliceKV(strings.Split(str, ";"))
	}
	return nil
}
