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

func (u *User) serialize(requester *User) *gabs.Container {
	res := gabs.New()
	_, _ = res.Set(u.ID, "id")
	_, _ = res.Set(u.Name, "name")
	if requester != nil && IsLoggedIn(requester.Role) {
		_, _ = res.Set(u.ProfileSettings, "profileSettings")
		_, _ = res.Set(u.Email, "email")
		_, _ = res.Set(u.Role, "role")

		_, _ = res.Set(u.Gender, "gender")
		_, _ = res.Set(u.EntryYear, "entry")
		_, _ = res.Set(u.ProfileCover, "profileCover")
		_, _ = res.Set(u.ProfileBoard, "profileBoard")
		if full || u.isClassPublic() {
			_, _ = res.Set(u.Class, "class")
		}
	}
	return res
}
