# Data

This document describes the data structure and relevant restrictions.

---

### 1. UserData
```go
type User struct {
    ID    uint16
    Email string
    Role  uint8
    // Personal information:
    Name      string
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
    UpdateDate uint64
    CreateDate uint64
}
```
#### Read operation
- Common information: id, profileSettings, profileCover, profileBoard
- Personal information: name, gender, entryYear, class, role; achievements and annual ranks
  <br>At least one requirement met to gain access:
  + Be the user himself
  + The profile is unlocked
    + However, with "class", there is an additional requirement is "class" being public
    + With "achievements", there is an additional requirement is "achievements" being public
    + With "annual ranks", there is an additional requirement is "annual ranks" being public
  + The requester is in the manager group
- Secret information: email, birthday, phone, updateDate, createDate
  <br>At least one requirement met to gain access:
  + Be the user himself
  + The requester is in the manager group
#### Write operation
- A user can edit his profile with the following acceptable fields: profileSettings, profileCover, profileBoard
- Managers can edit the role, achievements and annual ranks of other users
  + With class managers (Class Secretary, Class Deputy Secretary), there are additional restrictions:
    + The target user must be in the same class as the class manager
    + The target is ranked lower except Guest
    + The manager can only work with "Regular Member" and "Certified Member" roles
  + With global managers (Secretary, Deputy Secretary), there are additional restrictions:
    + The target is ranked lower except Guest
    + The manager can only work with "Regular Member", "Certified Member", "Class Secretary", "Class Deputy Secretary" roles
- UpdateDate, CreateData are handled by the system
- ID, Email, and the remaining are not overridable without special access
