# Viewer-Go

Viewer is a small application written in go. It is a web interface for browsing a given directory
of a server, which includes easy uploading and simple file management.

## Project Dependencies

Golang >=1.8  
gorilla/mux  
gorilla/sessions  
gorilla/securecookie  
crypto/bcrypt  
mattn/go-sqlite3

## Setup

- Get the code and dependencies:

```bash
go get github.com/gorilla/mux
go get github.com/gorilla/sessions
go get github.com/gorilla/securecookie
go get golang.org/x/crypto/bcrypt
go get github.com/mattn/go-sqlite3

go get github.com/ludw1gjj/viewer-go
```

- Go to viewer-go directory:

`cd {Your GOPATH}/src/github.com/ludw1gjj/viewer-go`

- To run in development:

`go run main.go -port=3000 -dbFile=viewer.db -configFile=config.json`

- Access via:

`http://localhost:3000`

- To build and run:  
  `go build`  
  `./viewer-go -port=3000 -dbFile=viewer.db -configFile=config.json`

## Notes

### Cross Compile

This will not cross compile with ease as the project is using mattn/go-sqlite3. If you want to save
the trouble, build on the platform you're targeting.

### Database - viewer.db

This project uses SQLite database. When the app is run it will check for viewer.db (or whatever you have inputted via the `-dbFile` flag, viewer.db is the default) file in the current directory, and
if not found, the file will automatically be created which includes the users table and a default
user.

The default user has these values:  
**Username:** admin  
**First Name:** John  
**Last Name:** Smith  
**Password:** password  
**Directory Root:** ./admin  
**IsAdmin:** true

As a result of creating the default user, a folder named "admin" will be created in the current directory. You may want
to delete this folder and change the Username, Directory Root and other values of this user.

### Configuration File - config.json

The file is a JSON file which contains two 32 byte length keys used for cookies. If not found on app startup, it will be created automatically in the current directory. The default file is config.json unless you used the `-configFile` flag to specify a different value.

Example config.json:

```bash
{
  "cookie": {
    "authorisation_key": "Akjp/SxHWvAB8e3MZ8T8qKoiICdQAer1UH7qwJTR4aw=",
    "encryption_key": "Ur6N8WGD4tW40semAKo5TusDHwaPZNzA9mUDHeP7EAA="
  }
}
```
