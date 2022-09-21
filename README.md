# QLDV

`QLDV` is the website of the youth union of Di An school (`DAS`).

## Project structure

- `migrator`: converts Excel files containing students' data to appropriate database format
- `server`: backend code
- `web`: frontend code

## Project documentation
- [API Specification](API.md)
- [Data Specification](Data.md)

## Project setup

### 1. Backend-side
**Environmental variables:**
```
GOOGLE_APPLICATION_CREDENTIALS=firebase.json
sql=user:pass@tcp(127.0.0.1:3306)/das?charset=utf8mb4&parseTime=True&loc=Local
```

## Project overview

### I. Structure
- `Youth Union` = Đoàn thanh niên
- `Class` = Chi đoàn

### II. Roles
- `Secretary` = Bí thư đoàn
  - Having full access to all features
- `Deputy Secretary` = Phó bí thư đoàn
  - Having full access to certain features
- `Class Secretary` = Bí thư chi đoàn
  - Having full access to all features related to the relevant class
- `Class Deputy Secretary` = Phó bí thư chi đoàn
  - Having full access to certain features related to the relevant class
- `Regular Member` = Thành viên (chưa gia nhập)
  - Having access to personal data and other (limited) features
- `Certified Member` = Đoàn viên
  - Having access to personal data, activities and other related features
- `Guest`
  - Read-only mode with limited access

In addition, there are shorthand terms to refer multiple roles:
- `User`: who has been logged in; in other words, are who belong to any roles except `Guest`
- `Member`: who belong to `Class`; in other words, are who in the group `Class Secretary`, `Class Deputy Secretary`, `Regular Member`, `Certified Member`
- `Class Manager`: who are in the following groups `Class Secretary`, `Class Deputy Secretary`
- `Global Manager`: who are in the following groups `Secretary`, `Deputy Secretary`
- `Manager`: who are either `Class Manager` or `Global Manager`

### III. Features

#### 1. Common pages
- Home page
  - Shows images and videos related to `DAS` and its youth union
  - Lists all clubs
  - Lists **happening** and **upcoming** events in the **current** month
    - `User` can join contest events if allowed and possible
- Committee page
  - Shows the structure of the committee and its members
- Event List page
  - Able to filter and sort events
  - Lists all events with respect to the settings
  - `User` can click a contest event to see further information if allowed
- Post List page
  - Able to view, filter and sort posts
  - Able to see post statistics
- Post page
  - Able to read the relevant post
  - `User` can like and see their reaction status
  - `User` who read the page for the first time will contribute one view to the counter

#### 2. Event pages

- Contest Overview page **[`User` only]**
  - Able to view general information about the contest
  - Able to view the participation history
  - Able to view the overall statistics
  - Can join the contest if allowed
- Contest page **[`User` only]**
  - Where people perform a contest

#### 3. Personal pages

- Profile page **[`User` only]**
  - The personal page of a `User`
- Class page **[`Member` only]**
  - Where `Member` belonging to the same group can discuss and share experience

#### 4. Management pages

Term `LCRUD`: List - Create - Read - Update - Delete

- Common Management page
  - Settings for the home page, the committee page
  - Other miscellaneous settings
- Post Management page
  - `LCRUD` for posts
  - View statistics
- Event Management page
  - `LCRUD` for events
  - View statistics
- Contest Management page
  - Able to `CRUD` a contest from its relevant event
  - View participation history and statistics
- User Management page
  - `LCRUD` for users
  - Can search and filter users
  - Can edit achievements and annual rankings of a user