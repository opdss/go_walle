package deploys

import (
	"encoding/json"
	"fmt"
)

type RemoteErrs map[int64]error

func (r RemoteErrs) Error() string {
	res := ""
	for k, v := range r {
		res = fmt.Sprintf("[%d]%s;%s", k, v, res)
	}
	return res
}

func (r RemoteErrs) String() string {
	res, _ := json.Marshal(r)
	return string(res)
}

func (r RemoteErrs) HasSuccess() bool {
	for _, v := range r {
		if v == nil {
			return true
		}
	}
	return false
}
