package security

const RoleGuest uint8 = 0
const RoleRegularMember uint8 = 1
const RoleCertifiedMember uint8 = 2
const RoleClassDeputySecretary uint8 = 3
const RoleClassSecretary uint8 = 4
const RoleDeputySecretary uint8 = 5
const RoleSecretary uint8 = 6
const RoleRoot uint8 = 7

const RoleGroupGuest uint8 = 0
const RoleGroupMember uint8 = 1
const RoleGroupClassManager uint8 = 2
const RoleGroupGlobalManager uint8 = 3
const RoleGroupRoot uint8 = 4

func GetRoleGroup(role uint8) uint8 {
	switch role {
	case RoleGuest:
		return RoleGroupGuest
	case RoleCertifiedMember, RoleRegularMember:
		return RoleGroupMember
	case RoleClassSecretary, RoleClassDeputySecretary:
		return RoleGroupClassManager
	case RoleSecretary, RoleDeputySecretary:
		return RoleGroupGlobalManager
	case RoleRoot:
		return RoleGroupRoot
	}
	return RoleGuest
}

func IsMember(role uint8) bool {
	return role == RoleRegularMember || role == RoleCertifiedMember
}

func IsManager(role uint8) bool {
	return GetRoleGroup(role) == RoleGroupClassManager || GetRoleGroup(role) == RoleGroupGlobalManager || GetRoleGroup(role) == RoleGroupRoot
}

func IsClassRole(role uint8) bool {
	return IsMember(role) || GetRoleGroup(role) == RoleGroupClassManager
}

func IsLoggedIn(role uint8) bool {
	return role != RoleGuest
}

func IsCertified(role uint8) bool {
	return role != RoleGuest && role != RoleRegularMember
}
