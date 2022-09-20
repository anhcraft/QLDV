package models

import (
	"das/security"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHasPrivilegeOver(t *testing.T) {
	user1 := &User{
		ID:   1,
		Role: security.RoleGuest,
	}
	user2 := &User{
		ID:   2,
		Role: security.RoleRegularMember,
	}
	user3 := &User{
		ID:   3,
		Role: security.RoleCertifiedMember,
	}
	user4 := &User{
		ID:   4,
		Role: security.RoleClassSecretary,
	}
	user5 := &User{
		ID:   5,
		Role: security.RoleSecretary,
	}
	user6 := &User{
		ID:   6,
		Role: security.RoleRoot,
	}
	assert.True(t, user1.HasPrivilegeOver(user1, 0), "user1 = user1; m = 0")
	assert.True(t, user1.HasPrivilegeOver(user1, 1), "user1 = user1; m = 1")
	assert.False(t, user1.HasPrivilegeOver(user2, 0), "user1 < user2; m = 0")
	assert.False(t, user1.HasPrivilegeOver(user2, 1), "user1 < user2; m = 1")
	assert.False(t, user1.HasPrivilegeOver(user4, 0), "user1 < user4; m = 0")
	assert.False(t, user1.HasPrivilegeOver(user4, 1), "user1 < user4; m = 1")
	assert.False(t, user2.HasPrivilegeOver(user3, 0), "user2 - user3; m = 0")
	assert.False(t, user2.HasPrivilegeOver(user3, 1), "user2 - user3; m = 1")
	assert.False(t, user4.HasPrivilegeOver(user5, 0), "user4 < user5; m = 0")
	assert.False(t, user4.HasPrivilegeOver(user5, 1), "user4 < user5; m = 1")
	assert.True(t, user5.HasPrivilegeOver(user4, 0), "user4 < user5; m = 0")
	assert.True(t, user5.HasPrivilegeOver(user4, 1), "user4 < user5; m = 1")
	assert.True(t, user5.HasPrivilegeOver(user3, 0), "user3 < user5; m = 0")
	assert.True(t, user5.HasPrivilegeOver(user3, 1), "user3 < user5; m = 1")
	assert.True(t, user5.HasPrivilegeOver(user2, 0), "user2 < user5; m = 0")
	assert.True(t, user5.HasPrivilegeOver(user2, 1), "user2 < user5; m = 1")
	assert.True(t, user5.HasPrivilegeOver(user1, 0), "user1 < user5; m = 0")
	assert.True(t, user5.HasPrivilegeOver(user1, 1), "user1 < user5; m = 1")
	assert.True(t, user4.HasPrivilegeOver(user3, 0), "user3 < user4; m = 0")
	assert.False(t, user4.HasPrivilegeOver(user3, 1), "user3 - user4; m = 1")
	assert.True(t, user4.HasPrivilegeOver(user2, 0), "user2 < user4; m = 0")
	assert.False(t, user4.HasPrivilegeOver(user2, 1), "user2 - user4; m = 1")
	assert.True(t, user4.HasPrivilegeOver(user1, 0), "user1 < user4; m = 0")
	assert.False(t, user4.HasPrivilegeOver(user1, 1), "user1 - user4; m = 1")
	assert.True(t, user6.HasPrivilegeOver(user5, 0), "user5 < user6; m = 0")
	assert.True(t, user6.HasPrivilegeOver(user5, 1), "user5 < user6; m = 1")
	assert.True(t, user6.HasPrivilegeOver(user2, 1), "user2 < user6; m = 0")
	assert.True(t, user6.HasPrivilegeOver(user2, 1), "user2 < user6; m = 0")
	assert.False(t, user1.HasPrivilegeOver(user6, 0), "user1 < user6; m = 0")
	assert.False(t, user1.HasPrivilegeOver(user6, 1), "user1 < user6; m = 1")
	assert.False(t, user5.HasPrivilegeOver(user6, 0), "user5 < user6; m = 0")
	assert.False(t, user5.HasPrivilegeOver(user6, 1), "user5 < user6; m = 1")
}
