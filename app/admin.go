package usrcanettocan

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

func getAdminPath() string {
	frontendPath, found := os.LookupEnv("ADMIN_PATH")

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

func getAdminPort() int {
	var adminPort int
	adminPort = 9402

	adminPortFromEnvStr, ok := os.LookupEnv("ADMIN_PORT")

	if ok {
		apiPortFromEnv, err := strconv.Atoi(adminPortFromEnvStr)

		if err != nil {
			adminPort = apiPortFromEnv
		}
	}

	return adminPort
}

func ServeAdmin() {
	router := mux.NewRouter()

	adminPath := getAdminPath()

	log.Info(fmt.Sprintf("Serving admin from %s", adminPath))

	fh := frontendHandler{staticPath: getAdminPath(), indexPath: "index.html"}
	router.PathPrefix("/").Handler(fh)
	log.Info("Initializing HTTP Admin server...")

	loggingHandler := LoggingHTTPHandler(fh)

	srv := &http.Server{
		Handler:      loggingHandler,
		Addr:         fmt.Sprintf("0.0.0.0:%d", getAdminPort()),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	srv.ListenAndServe()
}
