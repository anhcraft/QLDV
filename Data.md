# Data

This document describes the data structure and relevant restrictions.

---

### 1. User
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

#### Role constants
```go
const RoleGuest uint8 = 0
const RoleRegularMember uint8 = 1
const RoleCertifiedMember uint8 = 2
const RoleClassSecretary uint8 = 3
const RoleClassDeputySecretary uint8 = 4
const RoleSecretary uint8 = 5
const RoleDeputySecretary uint8 = 6
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

### 2. Achievement
```go
type Achievement struct {
  UserId uint16
  User   User
  Title  string
  Year   uint16
  // Date stuff:
  UpdateDate uint64
  CreateDate uint64
}
```

#### Read operation
- `title`, `year`, `updateDate` and `createDate` are readable as long as the associated user data is readable by the requester
- The remaining fields are hidden

#### Write operation
- `title`, `year` are writeable as long as the associated user data is writeable by the requester
- The remaining fields are system-generated

### 3. Annual rank
```go
type AnnualRank struct {
  UserId     uint16
  User       User
  Year       uint16
  Level      uint8
  UpdateDate uint64
  CreateDate uint64
}
```

#### Read operation
- `level`, `year`, `updateDate` and `createDate` are readable as long as the associated user data is readable by the requester
- The remaining fields are hidden

#### Write operation
- `level`, `year` are writeable as long as the associated user data is writeable by the requester
- The remaining fields are system-generated

### 4. Post
```go
type Post struct {
  ID         uint32
  Link       string
  Title      string
  Headline   string
  Content    string
  Hashtag    string
  Privacy    uint8
  LikeCount  uint
  ViewCount  uint
  UpdateDate uint64
  CreateDate uint64
}
```

#### Privacy level
- The `privacy` field determines the lowest level of role which is able to access the post
- For example, if the `privacy` of a post is set to `Regular Member`, any roles can see it except `Guest`

#### Read operation
- All fields are readable as long as the privacy requirement is met

#### Write operation
- `Title`, `Headline`, `Content`, `Hashtag`, `Privacy` are writeable as long as the requester belongs to the global manager group
- The remaining fields are system-generated
