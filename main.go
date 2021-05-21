package main

import (
	"URLMaintainer/config"
	"log"
	"net/http"

	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	err    error
)

type URLMapper map[string]string

func init() {
	logger, err = zap.NewDevelopment()
	if err != nil {
		log.Fatal("unable to initialize logger")
	}
}

func getDefaultMux() *http.ServeMux {

	mux := http.NewServeMux()
	logger.Info("creating a new router")
	mux.HandleFunc("/", Test)
	return mux
}

func main() {
	r := getDefaultMux()
	mapHandler := config.URLRedirectMW(r, "config.yaml")
	http.ListenAndServe(":8082", mapHandler)
}

// Test function to test
func Test(w http.ResponseWriter, r *http.Request) {
	logger.Info("received a test request. Its working!!")
}
