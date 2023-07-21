package field

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"golang.org/x/exp/constraints"
)

type Slices[T constraints.Ordered] []T

func (s *Slices[T]) Scan(value any) error {
	*s = make([]T, 0)
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("type error")
	}
	if len(bytes) < 3 {
		return nil
	}
	return json.Unmarshal(bytes, &s)
}

func (s Slices[T]) Value() (driver.Value, error) {
	return s.String(), nil
}

func (s Slices[T]) String() string {
	if len(s) == 0 {
		return ""
	}
	r, _ := json.Marshal(s)
	return string(r)
}
