package usrcanettocan

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

func getFrontendPath() string {
	frontendPath, found := os.LookupEnv("FRONTEND_PATH")

	if !found {
		wd, _ := os.Getwd()
		frontendPath = filepath.Join(filepath.Dir(wd), "admin", "build")
	}

	return frontendPath
}

type frontendHandler struct {
	staticPath string
	indexPath  string
}

func (h frontendHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// prepend the path with the path to the static directory
	path := r.URL.Path
	path = filepath.Join(h.staticPath, path)
	// check whether a file exists at the given path
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func ServeAdmin() {
	router := mux.NewRouter()

	fh := frontendHandler{staticPath: getFrontendPath(), indexPath: "index.html"}
	router.PathPrefix("/").Handler(fh)
	log.Info("Initializing HTTP Admin server...")

	loggingHandler := LoggingHTTPHandler(fh)

	srv := &http.Server{
		Handler:      loggingHandler,
		Addr:         "0.0.0.0:9402",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	srv.ListenAndServe()
}
