package main

const RoleGuest uint8 = 0
const RoleRegularMember uint8 = 1
const RoleCertifiedMember uint8 = 2
const RoleClassSecretary uint8 = 3
const RoleClassDeputySecretary uint8 = 4
const RoleSecretary uint8 = 5
const RoleDeputySecretary uint8 = 6

func IsMember(role uint8) bool {
	return role == RoleRegularMember || role == RoleCertifiedMember
}

func IsClassRole(role uint8) bool {
	return role == RoleRegularMember || role == RoleCertifiedMember || role == RoleClassSecretary || role == RoleClassDeputySecretary
}

func IsLoggedIn(role uint8) bool {
	return role != RoleGuest
}

func IsCertified(role uint8) bool {
	return role != RoleGuest && role != RoleRegularMember
}

// CheckPrivilegeInGroup Checks whether the given "role" has the same or higher rank than "required" in the same group.
// - Class group: Class Secretary > Class Deputy Secretary > Certified Member > Regular Member
// - Global group: Secretary > Deputy Secretary
// Note: mixed ranks in different groups are considered invalid.
func CheckPrivilegeInGroup(role uint8, required uint8) bool {
	if required == RoleRegularMember {
		return role == RoleRegularMember || role == RoleCertifiedMember || role == RoleClassSecretary || role == RoleClassDeputySecretary
	} else if required == RoleCertifiedMember {
		return role == RoleCertifiedMember || role == RoleClassSecretary || role == RoleClassDeputySecretary
	} else if required == RoleClassDeputySecretary {
		return role == RoleClassSecretary || role == RoleClassDeputySecretary
	} else if required == RoleDeputySecretary {
		return role == RoleSecretary || role == RoleDeputySecretary
	} else {
		return role == required
	}
}
