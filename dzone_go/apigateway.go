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

type ValReply struct {
	Data string
	Err  string
}
type ValReply2 struct {
	Data map[string]interface{}
	Err  string
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
				var result ValReply2
				err = c.Call("Server.Validate", data, &result)
				if err != nil {
					w.Write([]byte(err.Error()))
				} else {
					req.Header.Set("PID", result.Data["pid"].(string))
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
	u, _ := url.Parse("http://localhost:" + LoadConfiguration().ApiService.Port)
	apiProxy := httputil.NewSingleHostReverseProxy(u)

	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/api/{rest:.*}", Validate(handler(apiProxy)))

	fmt.Println("API Gateway is up and running...")
	log.Fatal(http.ListenAndServe(":80", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}
