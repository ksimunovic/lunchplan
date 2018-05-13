package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io/ioutil"
	"os"
	"encoding/json"
)

var config Config

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"database"`
	ApiService struct {
		Port string `json:"port"`
	} `json:"api_service"`
	HtmlService struct {
		Port string `json:"port"`
	} `json:"html_service"`
	UserService struct {
		Port string `json:"port"`
	} `json:"user_service"`
	MealService struct {
		Port string `json:"port"`
	} `json:"meal_service"`
	TagService struct {
		Port string `json:"port"`
	} `json:"tag_service"`
	CalendarService struct {
		Port string `json:"port"`
	} `json:"calendar_service"`
}

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

func HandlerWrap(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		h.ServeHTTP(w, req)
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/static/").Handler(HandlerWrap(http.StripPrefix("/static/", http.FileServer(http.Dir("HtmlService/static")))))

	for _, route := range routes {
		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	fmt.Println("HTML Service is up and running...")
	log.Fatal(http.ListenAndServe(":"+LoadConfiguration().HtmlService.Port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}
