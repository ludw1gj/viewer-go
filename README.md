# Viewer
Viewer is a small application written in go. It is a web interface for browsing a given directory of a server, which 
includes easy uploading and simple file management.

## Project Dependencies
Golang >=1.8  
gorilla/mux  
gorilla/sessions  
lib/pq  
crypto/bcrypt  

## Setup Notes:
* Get the code and dependencies:
```
go get github.com/gorilla/mux
go get github.com/gorilla/sessions
go get github.com/lib/pq
go get golang.org/x/crypto/bcrypt

git clone github.com/FriedPigeon/viewer-go
```
* Go to viewer-go directory: `cd viewer-go`
* To run in development: `go run main.go -dev=true`
* Access via `http://localhost:3000`
* Or to build: `go build` and to run: `./go-mathfever`

## Authors
* **FriedPigeon**

