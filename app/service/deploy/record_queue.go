package deploy

import (
	"errors"
	"sync"
)

var RecordQueueNone = errors.New("none")

type RecordQueue struct {
	mux   *sync.RWMutex
	none  bool
	queue []any
}

func (r *RecordQueue) Push(ele any) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	if r.none {
		return RecordQueueNone
	}
	r.queue = append(r.queue, ele)
	return nil
}

func (r *RecordQueue) Get(idx int) any {
	r.mux.RLock()
	defer r.mux.RUnlock()
	l := len(r.queue)
	if l > idx {
		return r.queue[idx]
	} else {
		if r.none {
			return RecordQueueNone
		}
	}
	return nil
}

func (r *RecordQueue) None() {
	r.mux.Lock()
	r.none = true
	r.mux.Unlock()
}

type RecordQueueRead struct {
	queue *RecordQueue
	idx   int
}

func (rr *RecordQueueRead) Next() any {
	v := rr.queue.Get(rr.idx)
	if v != nil {
		return v
	}
	return nil
}
