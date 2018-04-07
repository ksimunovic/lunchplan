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

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/rpc"
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
	UserService struct {
		Port string `json:"port"`
	} `json:"userservice"`
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

func Validate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {

				data := map[string]interface{}{
					"bearerToken": bearerToken[1],
				}

				c, err := rpc.Dial("tcp", "127.0.0.1:"+LoadConfiguration().UserService.Port)
				if err != nil {
					fmt.Println(err)
					return
				}

				var result []byte
				var profile map[string]string
				err = c.Call("Server.Validate", data, &result)
				if err != nil {
					w.Write([]byte(err.Error()))
				} else {
					_ = json.Unmarshal(result, &profile)
					req.Header.Set("PID", profile["pid"])
					next(w, req)
				}

			}
			//TODO: Implementacija za cookie based auth
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
	logPath := "development.log"

	openLogFile(logPath)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	u, _ := url.Parse("http://localhost:" + LoadConfiguration().ApiService.Port)
	apiProxy := httputil.NewSingleHostReverseProxy(u)

	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/api/{rest:.*}", Validate(handler(apiProxy)))

	fmt.Println("API Gateway is up and running...")
	log.Fatal(http.ListenAndServe(":80", logRequest(handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router))))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Println(err)
		}
		log.Println(string(requestDump))

		handler.ServeHTTP(w, r)
	})
}

func openLogFile(logfile string) {
	if logfile != "" {
		lf, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)

		if err != nil {
			log.Fatal("OpenLogfile: os.OpenFile:", err)
		}

		log.SetOutput(lf)
	}
}
