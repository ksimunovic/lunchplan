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
	"time"
)

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"database"`
	ApiService struct {
		Port string `json:"port"`
		Host string `json:"host"`
	} `json:"api_service"`
	HtmlService struct {
		Host string `json:"host"`
	} `json:"html_service"`
	UserService struct {
		Host string `json:"host"`
	} `json:"user_service"`
	MealService struct {
		Host string `json:"host"`
	} `json:"meal_service"`
	TagService struct {
		Host string `json:"host"`
	} `json:"tag_service"`
	CalendarService struct {
		Host string `json:"host"`
	} `json:"calendar_service"`
}

var config Config

func LoadConfiguration() Config {
	if (Config{}) != config {
		return config
	}
	response, err := http.Get("http://configservice:50000")
	if err != nil {
		fmt.Printf("%s; ", err)
		fmt.Println("Trying again in 5 seconds...")
		time.Sleep(5 * time.Second)
		return LoadConfiguration()
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

func ServiceCall(method string, serviceHost string) http.HandlerFunc {
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

		if req.RequestURI != "/login" {
			if len(data["sid"].(string)) != 12 && len(data["sid"].(string)) != 24 {
				fmt.Println("Invalid sid in request")
				return
			}
		}

		getVars := mux.Vars(req)
		for k, v := range getVars {
			data["get_"+k] = v
		}

		c, err := rpc.Dial("tcp", serviceHost)
		if err != nil {
			fmt.Println(err)
			return
		}

		var result []byte
		jsonData, _ := json.Marshal(data)
		err = c.Call("Server."+method, jsonData, &result)
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
