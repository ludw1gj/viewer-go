package route

import (
	"github.com/FriedPigeon/viewer-go/handler"
	"github.com/gorilla/mux"
)

func Load() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", handler.Redirect).Methods("GET")
	r.HandleFunc("/viewer/{path:.*}", handler.Viewer).Methods("GET")
	r.HandleFunc("/folder", handler.CreateFolder).Methods("POST")
	r.HandleFunc("/upload", handler.Upload).Methods("POST")
	r.HandleFunc("/delete", handler.Delete).Methods("POST")
	r.HandleFunc("/delete-all", handler.DeleteAll).Methods("POST")

	return r
}
