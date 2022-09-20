package models

import (
	"das/security"
	"github.com/Jeffail/gabs/v2"
)

const GenderMale = false
const GenderFemale = true

type User struct {
	ID    uint16 `gorm:"autoIncrement;primaryKey"`
	Email string `gorm:"unique;not null"`
	Role  uint8
	// Personal information:
	Name      string `gorm:"not null"`
	Gender    bool
	Birthday  uint64
	EntryYear uint16
	Phone     string
	Class     string
	// Profile stuff:
	ProfileCover    string
	ProfileBoard    string
	ProfileSettings uint8
	// Date stuff:
	UpdateDate uint64 `gorm:"autoUpdateTime:milli"`
	CreateDate uint64 `gorm:"autoCreateTime:milli"`
}

func (u *User) IsProfileLocked() bool {
	return (u.ProfileSettings & 1) != 0
}

func (u *User) IsClassPublic() bool {
	return ((u.ProfileSettings & 2) != 0) && !u.IsProfileLocked()
}

func (u *User) IsAchievementPublic() bool {
	return ((u.ProfileSettings & 4) != 0) && !u.IsProfileLocked()
}

func (u *User) IsAnnualRankPublic() bool {
	return ((u.ProfileSettings & 8) != 0) && !u.IsProfileLocked()
}

// HasPrivilegeOver Checks whether the current user has superior privilege to the other
// Groups:
// - Guest
// - Member: Regular Member, Certified Member
// - Class Managers: Class Secretary, Class Deputy Secretary
// - Global Managers: Secretary, Deputy Secretary, Root
// Rules:
// - The user gains himself privilege
// - Roles in the same group are not overlapped each other
// Mode:
// - 0 (Group-Distinct) Global Managers > Class Managers > Member > Guest (default)
// - 1 (Global-Distinct) Global Managers > Class Managers - Member - Guest
func (u *User) HasPrivilegeOver(who *User, mode uint8) bool {
	if u.ID == who.ID {
		return true
	}
	distinctGroupTest := security.GetRoleGroup(u.Role) > security.GetRoleGroup(who.Role)
	if mode == 1 {
		return distinctGroupTest && security.GetRoleGroup(u.Role) == security.RoleGroupGlobalManager
	} else {
		return distinctGroupTest
	}
}

// Serialize Serializes this user's data into a gabs container
// Information may be hidden due to privacy concerns.
func (u *User) Serialize(requester *User) *gabs.Container {
	res := gabs.New()
	_, _ = res.Set(u.ID, "id")
	_, _ = res.Set(u.IsProfileLocked(), "settings", "profileLocked")
	_, _ = res.Set(u.IsClassPublic(), "settings", "classPublic")
	_, _ = res.Set(u.IsAchievementPublic(), "settings", "achievementPublic")
	_, _ = res.Set(u.IsAnnualRankPublic(), "settings", "annualRankPublic")
	_, _ = res.Set(u.ProfileCover, "profileCover")
	_, _ = res.Set(u.ProfileBoard, "profileBoard")

	// Group-distinct check
	accessLocked := requester.HasPrivilegeOver(u, 0)
	// Personal information
	if accessLocked || !u.IsProfileLocked() {
		_, _ = res.Set(u.Name, "name")
		if u.Gender == GenderMale {
			_, _ = res.Set("male", "gender")
		} else if u.Gender == GenderFemale {
			_, _ = res.Set("female", "gender")
		} else {
			_, _ = res.Set("unknown", "gender")
		}
		_, _ = res.Set(u.EntryYear, "entryYear")
		if accessLocked || u.IsClassPublic() {
			_, _ = res.Set(u.Class, "class")
		}
		_, _ = res.Set(u.Role, "role")
	}
	// Secret information
	if accessLocked {
		_, _ = res.Set(u.Email, "email")
		_, _ = res.Set(u.Birthday, "birthday")
		_, _ = res.Set(u.Phone, "phone")
		_, _ = res.Set(u.UpdateDate, "updateDate")
		_, _ = res.Set(u.CreateDate, "createDate")
	}
	return res
}
