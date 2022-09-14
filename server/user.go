package main

import "github.com/Jeffail/gabs/v2"

const GenderMale = false
const GenderFemale = true

type User struct {
	ID    int    `gorm:"autoIncrement;primaryKey"`
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
}

func (u *User) isProfileLocked() bool {
	return (u.ProfileSettings & 1) == 1
}

func (u *User) isClassPublic() bool {
	return (u.ProfileSettings & 2) == 2
}

func (u *User) isAchievementPublic() bool {
	return (u.ProfileSettings & 4) == 4
}

func (u *User) isRatePublic() bool {
	return (u.ProfileSettings & 8) == 8
}

// hasPrivilegeOver Checks whether the current user has superior privilege to the other
// Groups:
// - Guest
// - Member: Regular Member, Certified Member
// - Class Managers: Class Secretary, Class Deputy Secretary
// - Global Managers: Secretary, Deputy Secretary
// Rules:
// - The user gains himself privilege
// - Roles in the same group are not overlapped each other
// Mode:
// - 0 (Group-Distinct) Global Managers > Class Managers > Member > Guest (default)
// - 1 (Class-Distinct) Global Managers > Class Managers - Member - Guest
func (u *User) hasPrivilegeOver(who *User, mode uint8) bool {
	if u.ID == who.ID {
		return true
	}
	distinctGroupTest := GetRoleGroup(u.Role)-GetRoleGroup(who.Role) > 0
	if mode == 1 {
		return distinctGroupTest && GetRoleGroup(u.Role) == RoleGroupGlobalManager
	} else {
		return distinctGroupTest
	}
}

// serialize Serializes this user's data into a gabs container
// Information may be hidden due to privacy concerns.
// - Common information: id, profileSettings, profileCover, profileBoard
// - Personal information: name, gender, entryYear, class, role
//    (*) At least one requirement met to gain access:
//      + Be the user himself
//      + The profile is unlocked
//        However, with "class", there is an additional requirement is "class" being public
//      + The requester is in the manager group
// - Secret information: email, birthday, phone
//    (*) At least one requirement met to gain access:
//      + Be the user himself
//      + The requester is in the manager group
func (u *User) serialize(requester *User) *gabs.Container {
	res := gabs.New()
	_, _ = res.Set(u.ID, "id")
	_, _ = res.Set(u.isProfileLocked(), "profile.settings.profileLocked")
	_, _ = res.Set(u.isClassPublic(), "profile.settings.classPublic")
	_, _ = res.Set(u.isAchievementPublic(), "profile.settings.achievementPublic")
	_, _ = res.Set(u.isRatePublic(), "profile.settings.ratePublic")
	_, _ = res.Set(u.ProfileCover, "profile.cover")
	_, _ = res.Set(u.ProfileBoard, "profile.board")

	// Group-distinct check
	accessLocked := requester.hasPrivilegeOver(u, 0)
	// Personal information
	if accessLocked || !u.isProfileLocked() {
		_, _ = res.Set(u.Name, "name")
		if u.Gender == GenderMale {
			_, _ = res.Set("male", "gender")
		} else {
			_, _ = res.Set("female", "gender")
		}
		_, _ = res.Set(u.EntryYear, "entryYear")
		if accessLocked || u.isClassPublic() {
			_, _ = res.Set(u.Class, "class")
		}
		_, _ = res.Set(u.Role, "role")
	}
	// Secret information
	if accessLocked {
		_, _ = res.Set(u.Email, "email")
		_, _ = res.Set(u.Birthday, "birthday")
		_, _ = res.Set(u.Phone, "phone")
	}
	return res
}
