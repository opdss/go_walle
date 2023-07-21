package ssh

import (
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
	res := &Envs{
		kvs: make(map[string]string),
		mux: &sync.RWMutex{},
	}
	for k, v := range source {
		res.kvs[k] = fmt.Sprintf("%v", v)
	}
	return res
}

func NewEnvsBySliceKV(source []string) *Envs {
	res := &Envs{
		kvs: make(map[string]string),
		mux: &sync.RWMutex{},
	}
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

func (e *Envs) Pick(keys ...string) *Envs {
	e.mux.RLock()
	defer e.mux.RUnlock()
	res := &Envs{
		kvs: make(map[string]string),
		mux: &sync.RWMutex{},
	}
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

func (e *Envs) Empty() bool {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return len(e.kvs) == 0
}
