package main

const RoleGuest = 0
const RoleRegularMember = 1
const RoleCertifiedMember = 2
const RoleClassSecretary = 3
const RoleClassDeputySecretary = 4
const RoleSecretary = 5
const RoleDeputySecretary = 6

func IsMember(role int) bool {
	return role == RoleRegularMember || role == RoleCertifiedMember || role == RoleClassSecretary || role == RoleClassDeputySecretary
}

func IsLoggedIn(role int) bool {
	return role != RoleGuest
}

// CheckPrivilegeInGroup Checks whether the given "role" has the same or higher rank than "required" in the same group.
// - Class group: Class Secretary > Class Deputy Secretary > Certified Member > Regular Member
// - Global group: Secretary > Deputy Secretary
// Note: mixed ranks in different groups are considered invalid.
func CheckPrivilegeInGroup(role int, required int) bool {
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
