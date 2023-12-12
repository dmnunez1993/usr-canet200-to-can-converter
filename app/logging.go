package usrcanettocan

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/felixge/httpsnoop"
	log "github.com/sirupsen/logrus"
)

const EmptySource = ""

func getDebugLvl() log.Level {
	logLevelStr, found := os.LookupEnv("LOG_LEVEL")
	if !found {
		return log.InfoLevel
	}

	if logLevelStr == "DEBUG" {
		return log.DebugLevel
	}

	if logLevelStr == "INFO" {
		return log.InfoLevel
	}

	if logLevelStr == "WARN" {
		return log.WarnLevel
	}

	if logLevelStr == "ERROR" {
		return log.ErrorLevel
	}

	return log.InfoLevel
}

func init() {
	log.SetOutput(os.Stdout)
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "02-01-2006 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
	log.SetLevel(getDebugLvl())
}

func getDetails() (string, int) {
	_, file, line, _ := runtime.Caller(2)
	return file, line
}

func LogDebug(source string, msg string) {
	file, line := getDetails()
	debugMessage := msg

	if source != EmptySource {
		debugMessage = fmt.Sprintf("%s: %s", source, msg)
	}

	log.Debugf("%s:%d %s", file, line, debugMessage)
}

type logrusHTTPHandler struct {
	handler http.Handler
}

type logrusResponseWriter struct {
	size int
	code int
	w    http.ResponseWriter
}

func (lw *logrusResponseWriter) Header() http.Header {
	return lw.w.Header()
}

func (lw *logrusResponseWriter) Write(data []byte) (int, error) {
	var err error
	var size int
	metrics := httpsnoop.CaptureMetricsFn(lw.w, func(w http.ResponseWriter) {
		size, err = w.Write(data)
	})
	lw.code = metrics.Code
	lw.size += size
	return size, err
}

func (lw *logrusResponseWriter) WriteHeader(statusCode int) {
	lw.w.WriteHeader(statusCode)
}

func makeLogrusHTTPWriter(w http.ResponseWriter) *logrusResponseWriter {
	writer := logrusResponseWriter{size: 0, w: w}
	return &writer
}

func (h logrusHTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	lw := makeLogrusHTTPWriter(w)
	h.handler.ServeHTTP(lw, req)
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		host = req.RemoteAddr
	}

	uri := req.RequestURI
	url := req.URL

	// Requests using the CONNECT method over HTTP/2.0 must use
	// the authority field (aka r.Host) to identify the target.
	// Refer: https://httpwg.github.io/specs/rfc7540.html#CONNECT
	if req.ProtoMajor == 2 && req.Method == "CONNECT" {
		uri = req.Host
	}
	if uri == "" {
		uri = url.RequestURI()
	}

	now := time.Now()

	//metrics := httpsnoop.CaptureMetrics(h.handler, lw, req)
	responseSize := lw.size

	method := req.Method
	message := fmt.Sprintf(
		"Request Received: %s - - [%s] \"%s %s\" %d %d",
		host,
		now.Format("02/Jan/2006:15:04:05 -0700"),
		method,
		uri,
		lw.code,
		responseSize,
	)
	log.Info(message)
}

func LoggingHTTPHandler(h http.Handler) http.Handler {
	return logrusHTTPHandler{handler: h}
}

type LogWriter struct {
}

func (lw LogWriter) Write(p []byte) (int, error) {
	log.Info(string(p))
	return len(p), nil
}

func NewLogWriter() LogWriter {
	return LogWriter{}
}
