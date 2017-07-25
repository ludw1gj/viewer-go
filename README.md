# Viewer
Viewer is a small application written in go. It is a web interface for browsing a given directory of a server, which 
includes easy uploading and simple file management.

## Project Dependencies
Golang >=1.8  
gorilla/mux  
gorilla/sessions  
gorilla/securecookie  
lib/pq  
mattn/go-sqlite3  
crypto/bcrypt  

## Setup
* Get the code and dependencies:
```
go get github.com/gorilla/mux
go get github.com/gorilla/sessions
go get github.com/gorilla/securecookie
go get github.com/lib/pq
go get github.com/mattn/go-sqlite3
go get golang.org/x/crypto/bcrypt

git clone github.com/FriedPigeon/viewer-go
```
* Go to viewer-go directory:  
`cd viewer-go`
* To run in development:  
`go run main.go -dev=true`
* Access via:  
`http://localhost:3000`
* To build and run:  
`go build`  
`./viewer-go`

## Notes
#### Static Files - /static/*  
A static file handler is included and can be used when invoking the dev flag set to true: `-dev=true`

#### Database - viewer.db
This project uses SQLite database. When the app is run it will check for a viewer.db, and if not found the file will 
automatically be created which includes the users table and a default user.  

The default user has these values:  
**Username:** admin  
**First Name:** John  
**Last Name:** Smith  
**Password:** password  
**Directory Root:** ./admin  
**IsAdmin:** true

As a result of creating the default user, a folder named "admin" will be created in the current directory. You may want
to delete this folder and change the Username, Directory Root and other values of this user.

#### Configuration File - config.json  
The file is a JSON file which contains two 32 byte length keys used for cookies.  

Example config.json:   
```
{
  "cookie": {
    "authorisation_key": "Akjp/SxHWvAB8e3MZ8T8qKoiICdQAer1UH7qwJTR4aw=",
    "encryption_key": "Ur6N8WGD4tW40semAKo5TusDHwaPZNzA9mUDHeP7EAA="
  }
}
```

## Authors
* **FriedPigeon**

