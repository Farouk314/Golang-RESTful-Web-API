# Go RESTful API

This is a sample RESTful API using Go.

# Normal execution

`go build`
`./gorestapi`

# Docker build

`docker gorestapi -t build .`
`docker run -rm -d -p 8000:8000/tcp gorestapi:latest`

# Running all tests
`go test ../...`

# Using godocs
`godoc -http=:6060 -play`

# Exposed Endpoints

Note: Send username (“userA/userB” in memory) in Auth header
Endpoints allowing to create, update and delete certificates:

/certificates/{id} – (“POST”) – Create certificate
Request body needs:
“title”required, “note”, and “createdAt(RFC3339 format)”required
as application/json

/certificates/{id} – (“PATCH”) – Update certificate
Request body may have:
“title” (cannot be empty), “note”
as application/json

/certificates/{id} – (“DELETE”) – Delete certificate
No body needed, userName taken from Auth.
/users/{userId}/certificates (“GET”) – Get user’s certificates
No body needed, userName taken from Auth

/certificates/{id}/transfers (“PATCH”) – Create transfer
“transfer” : { “to” } required
as application/json
/certificates/{id}/transfers (“PUT”) – Accept transfer
No body needed, user taken from Auth.
