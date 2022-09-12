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

func CheckRole(role int, expected int) bool {
	if expected == RoleRegularMember {
		return role == RoleRegularMember || role == RoleCertifiedMember || role == RoleClassSecretary || role == RoleClassDeputySecretary
	} else if expected == RoleCertifiedMember {
		return role == RoleCertifiedMember || role == RoleClassSecretary || role == RoleClassDeputySecretary
	} else if expected == RoleClassDeputySecretary {
		return role == RoleClassSecretary || role == RoleClassDeputySecretary
	} else if expected == RoleDeputySecretary {
		return role == RoleSecretary || role == RoleDeputySecretary
	} else {
		return role == expected
	}
}
