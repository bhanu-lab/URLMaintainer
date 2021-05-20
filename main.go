package main

import (
	"URLMaintainer/config"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	err    error
)

type URLMapper map[string]string

var PathToURLMap URLMapper

func init() {
	logger, err = zap.NewDevelopment()
	if err != nil {
		log.Fatal("unable to initialize logger")
	}
	PathToURLMap = make(map[string]string)
	PathToURLMap["mail"] = "https://mail.google.com/mail/u/0/#inbox"
}

func main() {
	r := mux.NewRouter()
	fmt.Println("creating a new router")
	r.Use(URLRedirectMW)
	fmt.Println("registerred middle ware")
	r.HandleFunc("/{app}", Test)
	http.ListenAndServe(":8082", r)
}

// URLRedirectMW middleware function which acts on request and redirects URL
// to another URL when it sees path matching another URL defined in yaml file
func URLRedirectMW(next http.Handler) http.Handler {

	// combining 2 different map configs from map built from yaml file and local map config
	PathToURLMap := FuseMap(PathToURLMap, config.YamlMapBuilder())
	logger.Info("Map fused")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("received a request")
		vars := mux.Vars(r)
		app := vars["app"]
		if redirectURL, ok := PathToURLMap[app]; ok {
			logger.Info("URL is redirecting to %s since received app is %s \n", zap.String("redirectURL", redirectURL), zap.String("app", app))
			http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
		} else {
			logger.Info("no mapping defined for app %s \n", zap.String("app", app))
		}

		next.ServeHTTP(w, r)
	})
}

func AppRouter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}

// Test function to test
func Test(w http.ResponseWriter, r *http.Request) {
	logger.Info("received a test request. Its working!!")
}

// FuseMap combines 2 maps. If both m1, m2 has same keys m2 is retained
func FuseMap(m1, m2 URLMapper) URLMapper {
	for key, value := range m1 {
		if _, ok := m2[key]; !ok {
			m2[key] = value
		}
	}
	fmt.Printf("Final combined map %+v \n", m2)
	return m2
}
