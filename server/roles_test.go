package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoleFunc(t *testing.T) {
	assert.False(t, IsMember(RoleGuest), "Guest is not a member")
	assert.False(t, IsMember(RoleSecretary), "Secretary is not a member")
	assert.True(t, IsMember(RoleRegularMember), "Regular Member is a member")
	assert.True(t, IsMember(RoleCertifiedMember), "Certified Member is a member")
	assert.False(t, IsLoggedIn(RoleGuest), "Guest must have not been logged in")
	assert.True(t, IsLoggedIn(RoleRegularMember), "Regular Member must have been logged in")
	assert.False(t, IsCertified(RoleGuest), "Guest is not certified")
	assert.False(t, IsCertified(RoleRegularMember), "Regular Member is not certified")
	assert.True(t, IsCertified(RoleCertifiedMember), "Certified Member is certified")
	assert.True(t, IsCertified(RoleSecretary), "Secretary is certified")
}

func TestRoleCheck(t *testing.T) {
	// role required true? T = 1; F = 0
	table := [][3]int{
		{RoleGuest, RoleGuest, 1},
		{RoleGuest, RoleRegularMember, 0},
		{RoleGuest, RoleClassSecretary, 0},
		{RoleGuest, RoleSecretary, 0},
		{RoleSecretary, RoleGuest, 0},
		{RoleRegularMember, RoleRegularMember, 1},
		{RoleClassSecretary, RoleRegularMember, 1},
		{RoleClassDeputySecretary, RoleRegularMember, 1},
		{RoleRegularMember, RoleDeputySecretary, 0},
		{RoleRegularMember, RoleSecretary, 0},
		{RoleGuest, RoleRegularMember, 0},
		{RoleSecretary, RoleRegularMember, 0},
		{RoleDeputySecretary, RoleRegularMember, 0},
		{RoleClassSecretary, RoleClassDeputySecretary, 1},
		{RoleSecretary, RoleDeputySecretary, 1},
		{RoleSecretary, RoleClassSecretary, 0},
		{RoleSecretary, RoleClassDeputySecretary, 0},
		{RoleDeputySecretary, RoleClassSecretary, 0},
		{RoleDeputySecretary, RoleClassDeputySecretary, 0},
		{RoleClassSecretary, RoleCertifiedMember, 1},
		{RoleClassDeputySecretary, RoleCertifiedMember, 1},
		{RoleCertifiedMember, RoleCertifiedMember, 1},
		{RoleRegularMember, RoleCertifiedMember, 0},
	}

	for _, v := range table {
		if v[2] == 1 {
			assert.Truef(t, CheckPrivilegeInGroup(v[0], v[1]), "%v is not equivalent or inherited from %v", v[0], v[1])
		} else if v[2] == 0 {
			assert.Falsef(t, CheckPrivilegeInGroup(v[0], v[1]), "%v is not equivalent or inherited from %v", v[0], v[1])
		}
	}
}
