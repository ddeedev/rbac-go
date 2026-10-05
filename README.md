# RBAC (Role-based access control) Example

This project desing to implement as gRPC backend base
Using Hexagonal Architecture along side with Module-style
Adapter from Hexagonal Architecture will be located at x/ dir with their related protobuf code

## Usage

### Docker

Simply run

```bash
docker compose up --build -d
```

## Local

For local usage user must have mognodb running and create file .env (follow .env.example).
After that user can run backend by

```bash
make serve
```

## gRPC

To use grpc service user must have `buf` installed  
Follow instruction here [installation](https://buf.build/docs/cli/installation/)
and grpcurl here [installation](https://github.com/fullstorydev/grpcurl)

gRPC server will run at port 50051 at default

user can list all available service by

```bash
grpcurl -plaintext localhost:50051 list

# result
# auth.AuthService
# grpc.reflection.v1.ServerReflection // ignore this one
# grpc.reflection.v1alpha.ServerReflection // ignore this one
# user.UserService
```

All available service is auth.AuthService and user.UserService

for auth Service endpoints

```bash
grpcurl -plaintext localhost:50051 list auth.AuthService
# result
# auth.AuthService.Login
# auth.AuthService.Register
```

for user Service endpoints

```bash
grpcurl -plaintext localhost:50051 list user.UserService

# result
# user.UserService.Create
# user.UserService.Delete
# user.UserService.Get
# user.UserService.List
# user.UserService.Update
```

### Use case

Registeration

Body:

```json
{ "name": "string", "email": "string", "password": "string" }
```

Response:

```json
{ "success" : bool }
```

```bash
grpcurl -plaintext \
  -d '{"name":"alice","email":"alice@example.com","password":"password123"}' \
  localhost:50051 auth.AuthService/Register

```

Login

Body:

```json
{ "username": "string", "password": "string" }
```

Response:

```json
{
  "accessToken": "string",
  "expiresIn": number,
  "user": {
    "id": "string",
    "name": "string",
    "email": "string"
  }
}
```

```bash
grpcurl -plaintext \
  -d '{"username":"ann@example.com","password":"password123"}' \
  localhost:50051 auth.AuthService/Login

```

Get User

```json
{ "id": "string" }
```

Response:

```json
{
  "id": "string",
  "name": "string",
  "email": "string",
  "createdAt": "string"
}
```

```bash
grpcurl -plaintext \
  -H "authorization: Bearer $TOKEN" \
  -d '{"id":"<user-id>"}' \
  localhost:50051 user.UserService/Get

```

List user
Body: None

Response:

```json
[{
  "id": "string",
  "name": "string",
  "email": "string",
  "createdAt": "string"
}, ...]
```

```bash
grpcurl -plaintext \
  -H "authorization: Bearer $TOKEN"  \
  localhost:50051 user.UserService/List
```

Update user
Body:

```json
{ "id": "string", "name": "string, "email": "string" }
```

Response:

```json
{
  "id": "string",
  "name": "string",
  "email": "string",
  "createdAt": "string"
}
```

```bash
grpcurl -plaintext \
  -H "authorization: Bearer $TOKEN" \
  -d '{"id":"<user-id>", "name":"alice_new", "email":"alice_new@example.com"}' \
  localhost:50051 user.UserService/Update

```

Delete user
Body:

```json
{ "id": "string" }
```

Response:

```json
{
  "id": "string",
  "name": "string",
  "email": "string",
  "createdAt": "string"
}
```

```bash
grpcurl -plaintext \
  -H "authorization: Bearer $TOKEN" \
  -d '{"id":"<user-id>"}' \
  localhost:50051 user.UserService/Delete

```

## User Management API

In order to use RestAPI gRCP has built-in grpc to rest transation user can visit at port 5001 or [swagger](http://localhost:5001/swagger) for better experience
path available are

Auth Service:

Method: Post
Path: <http://localhost:5001/v1/login>
Body:

```json
{
  "username": "string",
  "password": "string"
}
```

Response:

```json
{
  "accessToken": "string",
  "expiresIn": number,
  "user": {
    "id": "string",
    "name": "string",
    "email": "string"
  }
}
```

Method: Post
Path: <http://localhost:5001/v1/register>
Body:

```json
{
  "name": "string",
  "email": "string",
  "password": "string"
}
```

Response:

```json
{ "success" : bool }
```

User Service:

Method: Post
Header 'Authorization: Bearer $TOKEN'
Path: <http://localhost:5001/v1/users>

Body:

```json
{
  "name": "string",
  "email": "string",
  "password": "string"
}
```

Response:

```json
{
  "id": "string",
  "name": "string",
  "email": "string",
  "createdAt": "string"
}
```

Method: Get
Header 'Authorization: Bearer $TOKEN'
Path: <http://localhost:5001/v1/users/{id}>

Body:

Response:

```json
{
  "id": "string",
  "name": "string",
  "email": "string",
  "createdAt": "string"
}
```

Method: Get
Header 'Authorization: Bearer $TOKEN'
Path: <http://localhost:5001/v1/users>

Body:

Response:

```json
[{
  "id": "string",
  "name": "string",
  "email": "string",
  "createdAt": "string"
}, ...]
```

Method: Put
Header 'Authorization: Bearer $TOKEN'
Path: <http://localhost:5001/v1/users/{id}>

Body:

```json
{
  "name": "string",
  "email": "string"
}
```

Response:

```json
{
  "id": "string",
  "name": "string",
  "email": "string",
  "createdAt": "string"
}
```

Method: Delete
Header 'Authorization: Bearer $TOKEN'
Path: <http://localhost:5001/v1/users/{id}>

Body:

Response:

```json
{
  "id": "string",
  "name": "string",
  "email": "string",
  "createdAt": "string"
}
```
