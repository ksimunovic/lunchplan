package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/couchbase/gocb"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Session struct {
	Type string `json:"type,omitempty"`
	Pid  string `json:"pid,omitempty"`
}

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"database"`
	ApiGateway struct {
		Port string `json:"port"`
	} `json:"api_gateway"`
	ApiService struct {
		Port string `json:"port"`
	} `json:"apiservice"`
}

var config Config

func LoadConfiguration() Config {

	if (Config{}) != config {
		return config
	}

	response, err := http.Get("http://localhost:50000/")
	if err != nil {
		fmt.Printf("%s", err)
		return Config{}
	} else {
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		config := Config{}
		jsonErr := json.Unmarshal(body, &config)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		return config
	}
}

var bucket *gocb.Bucket

func Validate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				var session Session
				_, err := bucket.Get(bearerToken[1], &session)
				if err != nil {
					w.WriteHeader(401)
					w.Write([]byte(err.Error()))
					return
				}
				req.Header.Set("PID", session.Pid)
				bucket.Touch(bearerToken[1], 0, 3600)
				next(w, req)
			}
		} else {
			req.Header.Set("PID", "")
			next(w, req)
		}
	})
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = mux.Vars(r)["rest"]
		p.ServeHTTP(w, r)
	}
}

func main() {
	fmt.Println("Starting the Go server...")

	u, _ := url.Parse("http://localhost:" + LoadConfiguration().ApiService.Port)
	apiProxy := httputil.NewSingleHostReverseProxy(u)

	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/api/{rest:.*}", Validate(handler(apiProxy)))

	cluster, _ := gocb.Connect(LoadConfiguration().Database.Host)
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: LoadConfiguration().Database.Username,
		Password: LoadConfiguration().Database.Password,
	})
	bucket, _ = cluster.OpenBucket("default", "")
	log.Fatal(http.ListenAndServe(":80", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}
