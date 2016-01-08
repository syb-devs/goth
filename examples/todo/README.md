# Goth REST API example

This is an example app using Goth for a RESTful API.

Example using cURL for creating, retrieving, updating and deleting a resource:

```shell
curl -H "Content-Type: application/json" -X POST -d '{"username": "john.doe@domain.com","password": "123456"}'
 http://localhost:3000/users

curl -H "Content-Type: application/json" -X PUT -d '{"firstName": "John", lastName: "Doe"}' http://localhost:3000/users/568f8cf6e5c7a6026ae87670

curl http://localhost:3000/users/568f8cf6e5c7a6026ae87670

curl -X DELETE http://localhost:3000/users/568f8cf6e5c7a6026ae87670
```
