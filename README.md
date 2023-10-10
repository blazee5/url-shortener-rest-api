# url-shortener-rest-api
REST API client written in golang for url shortener app
## How to use
```
git clone https://github.com/blazee5/url-shortener-rest-api
go mod tidy
go run cmd/url-shortener/main.go
```
## Routes
```go
r.Post("/url", save.New(log, dao)) // Create URL
r.Get("/{alias}", redirect.New(log, dao)) // Get redirect to URL
r.Delete("/{alias}", delete.Delete(log, dao)) // Delete URL
```
## Technologies
Project is created with:
* <img height="25" width="25" src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/go/go-original.svg" />  Go (Golang)
* <img height="25" width="25" src="https://cdn.simpleicons.org/mongodb/#47A248" /> MongoDB
## Features
What is realised in this project:
* REST API client
* Functional and unit tests
* Logger (Slog)
* Graceful shutdown
