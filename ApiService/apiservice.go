package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/rpc"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"database"`
	UserService struct {
		Port string `json:"port"`
	} `json:"UserService"`
	ApiService struct {
		Port string `json:"port"`
	} `json:"ApiService"`
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

func UserService(method string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		var data map[string]interface{}
		_ = json.NewDecoder(req.Body).Decode(&data)
		if len(data) == 0 {
			data = map[string]interface{}{
				"sid": req.Header.Get("sid"),
			}
		} else {
			data["sid"] = req.Header.Get("sid")
		}

		c, err := rpc.Dial("tcp", "127.0.0.1:"+LoadConfiguration().UserService.Port)
		if err != nil {
			fmt.Println(err)
			return
		}

		var result []byte
		err = c.Call("Server."+method, data, &result)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write(result)
		}
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	fmt.Println("API Service is up and running...")
	log.Fatal(http.ListenAndServe(":"+LoadConfiguration().ApiService.Port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}