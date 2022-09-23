package security

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoleFunc(t *testing.T) {
	assert.False(t, IsClassRole(RoleGuest), "Guest is not a class role")
	assert.False(t, IsClassRole(RoleSecretary), "Secretary is not a class role")
	assert.True(t, IsClassRole(RoleRegularMember), "Regular Member is a class role")
	assert.True(t, IsClassRole(RoleCertifiedMember), "Certified Member is a class role")
	assert.False(t, IsLoggedIn(RoleGuest), "Guest must have not been logged in")
	assert.True(t, IsLoggedIn(RoleRegularMember), "Regular Member must have been logged in")
	assert.False(t, IsCertified(RoleGuest), "Guest is not certified")
	assert.False(t, IsCertified(RoleRegularMember), "Regular Member is not certified")
	assert.True(t, IsCertified(RoleCertifiedMember), "Certified Member is certified")
	assert.True(t, IsCertified(RoleSecretary), "Secretary is certified")
	assert.True(t, IsMember(RoleRegularMember), "Regular Member is a member")
	assert.False(t, IsMember(RoleClassSecretary), "Class Secretary is not a member")
}

func TestRoleGroup(t *testing.T) {
	// role group
	table := [][2]uint8{
		{RoleGuest, RoleGroupGuest},
		{RoleCertifiedMember, RoleGroupMember},
		{RoleRegularMember, RoleGroupMember},
		{RoleClassSecretary, RoleGroupClassManager},
		{RoleClassDeputySecretary, RoleGroupClassManager},
		{RoleSecretary, RoleGroupGlobalManager},
		{RoleDeputySecretary, RoleGroupGlobalManager},
	}

	for _, v := range table {
		assert.Equalf(t, v[1], GetRoleGroup(v[0]), "%v belongs to group %v", v[0], v[1])
	}
}
