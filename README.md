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
  -d '{"username":"alice@example.com","password":"password123"}' \
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
{ "id": "string", "name": "string", "email": "string" }
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

## JWT Guide

I chose "github.com/golang-jwt/jwt/v5" as main dependencies for JWT adapter
For this POC there is only one type of JWT, which is accessToken and 24h of TTL
To make it more secure I force JWT secret must not less than 32 letter
and use that secret to signed with method HMAC-SHA256 (HS256)

In order to verify JWT golang-jwt has method ParseWithClaims for convert string token
also WithValidMethods and WithExpirationRequired to verify user string token

In order to use JWT, I have auth_middleware to work with grpc endpoints handler

There is two endpoints that user able to access w/o JWT, which are Register and Login

And how to use that you can see above in part of gRPC and RestAPI

## Design Decisions & Assumptions

As recommend and I have some experience with Architecture — Hexagonal (ports & adapters).
Most of Business logic live in core/ (domain, ports and service), but adapters I decided it to located and module path (x/)
It make codebase more clean and easy to manage with grpc backend

**Port (interfaces)**

- ports.UserRepository
- ports.TokenManager
- ports.PasswordHasher

**Adapter**

- MongoDB (`x/user/user_repo.go`), if there is more collection it will located at `x/new/new_repo.go`
- GRPC (x/user/user_handler.go)
- JWT/HS256 (`x/auth/jwt`)

**Utils**

- Argon2id (`utils/crypto`)

**Transport.**

We only have to impelement gRPC server and then RestAPI will be grpc-gateway that transcodes HTTP↔gRPC,
so there is only one set of validation and auth rules for both of grpc and rest,
Allow grpc reflection to make grpcurl and swagger able to public access.

**Security**

- **Argon2id** for password hashing (memory-hard; preferred over bcrypt).
  Minimum password length is 8.
- **HS256 JWT** with a secret of ≥32 bytes and a positive TTL, enforced at
  construction.
- **No user enumeration:** an unknown email and a wrong password both return
  the same `invalid credentials` error.
- Passwords are never returned in any response (`json:"-"` on the domain model).
