package usrcanettocan

import (
	"encoding/json"
	"fmt"
	logging "log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	log "github.com/sirupsen/logrus"
)

const errorStatus = "error"
const successStatus = "success"

type errorResponse struct {
	Status string `json:"status"`
	Detail string `json:"detail"`
}

type configSuccessResponse struct {
	Status string                  `json:"status"`
	Config UsrCanetConverterConfig `json:"data"`
}

func writeError(statusCode int, err error, res http.ResponseWriter, req *http.Request) {
	data := errorResponse{
		Status: errorStatus,
		Detail: err.Error(),
	}

	log.Error(err.Error())
	res.WriteHeader(statusCode)
	json.NewEncoder(res).Encode(data)
}

func writeConfigSuccess(config UsrCanetConverterConfig, res http.ResponseWriter, req *http.Request) {
	data := configSuccessResponse{
		Status: successStatus,
		Config: config,
	}

	json.NewEncoder(res).Encode(data)
}

func getConverterConfig(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	config, err := GetOrCreateConverterConfig()

	if err != nil {
		writeError(http.StatusInternalServerError, err, res, req)
		return
	}

	writeConfigSuccess(config, res, req)
}

func setConverterConfig(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(req.Body)
	newConfig := UsrCanetConverterConfig{}
	err := decoder.Decode(&newConfig)

	if err != nil {
		writeError(http.StatusBadRequest, err, res, req)
		return
	}

	err = newConfig.Save()

	if err != nil {
		writeError(http.StatusInternalServerError, err, res, req)
		return
	}

	writeConfigSuccess(newConfig, res, req)
}

func setRestApiRoutes(r *mux.Router) {
	r.HandleFunc("/get_config/", http.HandlerFunc(getConverterConfig)).Methods("GET")
	r.HandleFunc("/set_config/", http.HandlerFunc(setConverterConfig)).Methods("POST")
}

func getApiPort() int {
	var apiPort int
	apiPort = 9401

	apiPortFronEnvStr, ok := os.LookupEnv("API_PORT")

	if ok {
		apiPortFromEnv, err := strconv.Atoi(apiPortFronEnvStr)

		if err != nil {
			apiPort = apiPortFromEnv
		}
	}

	return apiPort
}

func ServeRestApi() {
	r := mux.NewRouter()

	setRestApiRoutes(r.PathPrefix("/").Subrouter())

	c := cors.New(cors.Options{
		// Allow all origins since the app is to be runned locally
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowCredentials: true,
		Debug:            getDebugLvl() == log.DebugLevel,
	})
	c.Log = logging.New(NewLogWriter(), "[cors]", logging.LstdFlags)

	log.Info("Initializing HTTP API server...")
	handler := c.Handler(r)

	loggingHandler := LoggingHTTPHandler(handler)

	srv := &http.Server{
		Handler:      loggingHandler,
		Addr:         fmt.Sprintf("0.0.0.0:%d", getApiPort()),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	srv.ListenAndServe()
}
