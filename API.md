# API

This document explains all endpoints under the public API.

### Status code
- 200: on success
- 400: on failure

### General response format
Every response will be based on the following format:
```json
{
    "result": {},
    "error": "",
    "success": true
}
```
- "success": always present, to represent the status of the response:
    + If true, the data may be present and there must be no error
    + If false, there must be an error occurred
- "error": points out the exact error code
    + This may be empty in case of success status
- "result": the result returned for the requester
    + It can be empty no matter if the request is success or not

### "Optional authentication"
- This means a request may have "access-token" header field to unlock further access to an endpoint.
- In case of absence, the endpoint still **have to** accept the request and process under the guest (dummy) user.

### "Authentication required"
- This means a request must have "access-token" header field to unlock further access to an endpoint.

---

### GET /user/:id?
- Get a user's profile, achievements and annual ranks
- Optional authentication
- Queries:
  + profile = true/false
  + achievements = true/false
  + annual-ranks = true/false
- Note:
  + If the id is absent, this API acts like "getting your own profile". The "access-token" field is required to obtain the corresponding user's id, and in case of absence, the API returns error.
  + If the id is present, this API acts like "getting someone profile" which "someone" can be "self" or "another"
- Example response:
```json
{
  "result": {
    "profile": {
      "id": 120,
      "settings": {
        "classPublic": false
      },
      "entryYear": 2020
    },
    "achievements": [
      {
        "title": "",
        "year": 2021
      }
    ],
    "annual-ranks": [
      {
        "year": 2020,
        "level": 0
      }
    ]
  },
  "success": true
}
```

### POST /user/:id?
- Change a user's data such as information, achievements, annual ranks, etc
- Authentication required
- Example request:
```json
{
  "profile": {
    "settings": {
      "classPublic": false
    },
    "entryYear": 2020
  },
  "achievements": [
    {
      "title": "",
      "year": 2021
    }
  ],
  "annual-ranks": [
    {
      "rank": 0,
      "year": 2020
    }
  ],
}
```
- Note:
  + Only fields which were specified in the request are considered. However, not all of them are editable because of various reasons mentioned above. Some profile fields are modifiable if and only if certain requirements met. For example, the requester must be the user himself or has special permissions. In addition, a few fields are open to the managers only and members are prohibited to edit them. Besides, there are fields reserved and under read-only mode such as ID, Email, etc **(See the Data document for further information)**
  + "achievements" and "annual-ranks" must be present with empty values to reset the corresponding field; otherwise, there is no effect if they are absent in the response.
- Example response:
```json
{
  "result": {
    "profile": {
      "settings": true,
      "profileBoard": true,
      "role": true
    },
    "achievements": true,
    "annualRanks": true
  },
  "success": true
}
```

### GET /users/
- Lists and filters users
- Authentication required
- Example request:
```json
{
  "limit": 50,
  "filterName": "",
  "filterClass": "",
  "filterEmail": "",
  "filterRole": 1,
  "belowId": 0
}
```
- Note:
  + Only managers can do this request
    + Class managers are prohibited from fetching users out of their classes
    + Global managers can fetch all users
  + The system will determine which fields are included in the response, which is the same as `GET /user/:id?`
- Example response:

```json
{
  "result": [
    {
      "id": 0
    },
    {
      "id": 1
    }
  ],
  "success": true
}
```

### GET /user-stats/
- Fetches user statistics
- Authentication required
- Note:
  + Only global managers can do this request
- Example response:
```json
{
  "result": {
    "user-count-by-role": {
      "regular-member": 10,
      "certified-member": 10,
      "class-deputy-secretary": 10,
      "class-secretary": 10,
      "deputy-secretary": 10,
      "secretary": 10
    },
    "user-count-by-grade": {
      "grade-10": 40,
      "grade-11": 50,
      "grade-12": 30
    }
  },
  "success": true
}
```

### GET /user-profile-cover/
- Uploads and sets the specified image as profile cover
- Authentication required
- Example response:
```json
{
  "result": {
    "name": "file_name.png"
  },
  "success": true
}
```

---

### GET /post/:id
- Gets the given post's data
- Optional authentication
- Note:
  + When the requester is logged in, there will be additional fields related to post statistics: `viewed`, `liked`
  + An error may occur if the post is hidden due to privacy settings
- Example response:
```json
{
  "result": {
    "id": 1,
    "title": "Test post",
    "viewed": true
  },
  "success": true
}
```

### POST /post/:id?
- Updates or creates a post
- Authentication required
- Note:
  + For now, only global managers can make request to this endpoint
  + When the `id` param is absent, this means "creating"; otherwise it means "updating"
- Example request:
```json
{
  "title": "Test post",
  "content": "Hello world"
}
```
- Example response:
```json
{
  "result": {
    "id": 1
  },
  "success": true
}
```

### DELETE /post/:id
- Deletes a post
- Authentication required
- Note:
  + For now, only global managers can make request to this endpoint
- Example response:
```json
{
  "success": true
}
```

### GET /posts/
- Lists and filters posts
- Authentication required
- Example request:
```json
{
  "limit": 50,
  "filterHashtags": [
    "news"
  ],
  "belowId": 0,
  "sortBy": "like",
  "lowerThan": 100
}
```
- Note:
  + Certain posts may be hidden due to privacy settings
  + The system will determine which fields are included in the response, which is the same as `GET /post/:id`
- Example response:
```json
{
  "result": [
    {
      "id": 0
    },
    {
      "id": 1
    }
  ],
  "success": true
}
```

### POST /post-like/:id
- Switches
- Authentication required
- Supported statistic types:
  + `like` - the state is switched for e
- Example response:
```json
{
  "result": [
    {
      "id": 0
    },
    {
      "id": 1
    }
  ],
  "success": true
}
```

### POST /post-attachment/:id

### GET /post-hashtags/
