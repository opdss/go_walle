package constants

const StatusEnable = 1
const StatusDisable = 2

const SuperUserId = 1

func IsSuperUser(userId int64) bool {
	return userId == SuperUserId
}

type Role string

const RoleSuper Role = "super"
const RoleOwner Role = "owner"
const RoleMaster Role = "master"
const RoleDeveloper Role = "developer"

var roleLevel = map[Role]int{
	RoleDeveloper: 1,
	RoleMaster:    1 << 2,
	RoleOwner:     1 << 3,
	RoleSuper:     1 << 4,
}

func (r Role) Level() int {
	if v, ok := roleLevel[r]; ok {
		return v
	}
	return 0
}
