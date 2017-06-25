package route

import (
	"flag"
	"net/http"

	"github.com/FriedPigeon/viewer-go/handler"
)

func Load() {
	http.HandleFunc("/", handler.Redirect)
	http.HandleFunc("/viewer/", handler.Viewer)
	http.HandleFunc("/folder", handler.CreateFolder)
	http.HandleFunc("/upload", handler.Upload)
	http.HandleFunc("/delete", handler.Delete)
	http.HandleFunc("/delete-all", handler.DeleteAll)

	// static files handler in dev mode
	dev := flag.Bool("dev", false, "Use in development")
	flag.Parse()
	if *dev {
		fs := http.FileServer(http.Dir("./static"))
		http.Handle("/static/", http.StripPrefix("/static", fs))
	}
}
