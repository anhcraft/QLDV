package main

const RoleGuest = 0
const RoleRegularMember = 1
const RoleClassSecretary = 2
const RoleClassDeputySecretary = 3
const RoleSecretary = 4
const RoleDeputySecretary = 5

func IsMember(role int) bool {
	return role == RoleRegularMember || role == RoleClassSecretary || role == RoleClassDeputySecretary
}

func IsLoggedIn(role int) bool {
	return role != RoleGuest
}

func CheckRole(role int, expected int) bool {
	if expected == RoleRegularMember {
		return role == RoleRegularMember || role == RoleClassSecretary || role == RoleClassDeputySecretary
	} else if expected == RoleClassDeputySecretary {
		return role == RoleClassSecretary || role == RoleClassDeputySecretary
	} else if expected == RoleDeputySecretary {
		return role == RoleSecretary || role == RoleDeputySecretary
	} else {
		return role == expected
	}
}
