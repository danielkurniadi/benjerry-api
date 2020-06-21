# User API Schema

This documenation includes all API endpoints for user signup and registration. The purpose of having
the user API is for the purpose of prototyping such that session based authentication can be demonstrated in this project.

> Notes:
> - Timestamp used is in seconds. (i.e. need to times with 1000 in JavaScript)
> - Every failed request is expected to be responded with an `error_message` field. For example:
```json
{
  "Message": "Insufficient permissions"
}
```

## Login User

`POST api/users/login`

### Request 

#### Basic Auth Header:
| Name                  | Value                 | Description
| -----------------     | --------              | -----------
| `username`            | `String`              | username is a unique identifier for user
| `password`            | `String`              | raw user password


### Response 

#### Cookie:
| Name                  | Value                 | Description
| -----------------     | --------              | -----------
| `session_token`       | `String`              | UUID V4 session token for authentication

#### Body:

##### No Error
`HTTP 200 OK`
```
{ "Message": "Login Successfull" }
```
##### Error
`HTTP 401 Unauthorized`
```
{ "Message": "Login Successfull" }
```
---

## Register Member User

`POST api/users/signup`

### Request 

#### Basic Auth Header:
| Name                  | Value                 | Description
| -----------------     | --------              | -----------
| `username`            | `String`              | username to register user
| `password`            | `String`              | raw user password to register user


### Response 

#### Body:

##### No Error
`HTTP 201 CREATED`
```
{ "Message": "Account created successfully" }
```
##### Error
`HTTP 409 Conflict`
```
{ "Message": "User with existing username" }
```

---

## Register Admin User

`POST api/users/admin`

### Request 

#### Basic Auth Header:
| Name                  | Value                 | Description
| -----------------     | --------              | -----------
| `username`            | `String`              | username to register admin
| `password`            | `String`              | raw user password to register admin


### Response 

#### Body:

##### No Error
`HTTP 201 CREATED`
```
{ "Message": "Account created successfully" }
```
##### Error
`HTTP 409 Conflict`
```
{ "Message": "User with existing username" }
```
