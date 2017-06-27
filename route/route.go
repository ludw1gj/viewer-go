package route

import (
	"flag"
	"net/http"

	"github.com/FriedPigeon/viewer-go/handler"
	"github.com/gorilla/mux"
)

func Load() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", handler.Redirect).Methods("GET")
	r.HandleFunc("/viewer/{path:.*}", handler.Viewer).Methods("GET")
	r.HandleFunc("/about", handler.About).Methods("GET")

	r.HandleFunc("/upload", handler.Upload).Methods("POST")
	r.HandleFunc("/folder", handler.CreateFolder).Methods("POST")
	r.HandleFunc("/delete", handler.Delete).Methods("POST")
	r.HandleFunc("/delete-all", handler.DeleteAll).Methods("POST")
	r.NotFoundHandler = http.HandlerFunc(handler.NotFound)

	// static file handler in dev mode
	boolPtr := flag.Bool("dev", false, "Use in development")
	flag.Parse()
	if *boolPtr {
		fs := http.FileServer(http.Dir("./static"))
		r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	}

	return r
}
