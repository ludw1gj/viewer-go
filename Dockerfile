FROM golang:1.9.3-stretch

# Get dependencies
RUN go get github.com/gorilla/mux; \
  go get github.com/gorilla/sessions; \
  go get github.com/gorilla/securecookie; \
  go get golang.org/x/crypto/bcrypt; \
  go get github.com/mattn/go-sqlite3

# Copy files and install
COPY . /go/src/github.com/ludw1gj/viewer-go
RUN go install github.com/ludw1gj/viewer-go

# Set working directory so relative go templates work properly
WORKDIR /go/src/github.com/ludw1gj/viewer-go

# Run app and expose port
ENTRYPOINT /go/bin/viewer-go
EXPOSE 3000
