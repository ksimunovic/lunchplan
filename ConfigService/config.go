package main

import (
	"net/http"
	"log"
	"fmt"
	"encoding/json"
)

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"database"`
	Cookie64 string `json:"cookie_32"`
	Cookie32 string `json:"cookie_64"`
	ApiGateway struct {
		Port string `json:"port"`
	} `json:"api_gateway"`
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

func main() {

	config := &Config{
		Database: struct {
			Host     string `json:"host"`
			Username string `json:"username"`
			Password string `json:"password"`
		}{Host: "couchbase://localhost", Username: "Administrator", Password: "Administrator"},
		Cookie64:    "f56931b180a522c21c17ba9274d07b0934a366bef09a7b43db9006a06b97bc3cf4a4e2c2ae77d11d3c38c6f17310eac0fbbb8da5eb1d1ae57eaef133abbf5ec8",
		Cookie32:    "20c2f69c5d03cd82deb2203bc9b1701fad3b7c3a7d3d045c43bb6994ccb44afd",
		ApiGateway:  struct{ Port string `json:"port"` }{Port: "50000"},
		ApiService:  struct{ Port string `json:"port"` }{Port: "50001"},
		HtmlService:  struct{ Port string `json:"port"` }{Port: "50002"},
		UserService: struct{ Port string `json:"port"` }{Port: "50003"},
		MealService: struct{ Port string `json:"port"` }{Port: "50004"},
		TagService: struct{ Port string `json:"port"` }{Port: "50005"},
		CalendarService: struct{ Port string `json:"port"` }{Port: "50006"},
	}

	configJson, _ := json.Marshal(config)

	fmt.Println("Config.json is up and running...")
	log.Fatal(http.ListenAndServe(":50000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		w.Write(configJson)
	})))
}
