package field

import "go-walle/app/internal/constants"

const StatusEnable = 1
const StatusDisable = 2

type Status uint8

func (f Status) IsEnable() bool {
	return f == constants.StatusEnable
}

func (f Status) IsDisable() bool {
	return f == constants.StatusDisable
}
