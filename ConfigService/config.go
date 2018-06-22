package main

import (
	"net/http"
	"log"
	"encoding/json"
	"os"
	"net"
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
		Host string `json:"host"`
	} `json:"api_gateway"`
	ApiService struct {
		Port string `json:"port"`
		Host string `json:"host"`
	} `json:"api_service"`
	HtmlService struct {
		Port string `json:"port"`
		Host string `json:"host"`
	} `json:"html_service"`
	UserService struct {
		Port string `json:"port"`
		Host string `json:"host"`
	} `json:"user_service"`
	MealService struct {
		Port string `json:"port"`
		Host string `json:"host"`
	} `json:"meal_service"`
	TagService struct {
		Port string `json:"port"`
		Host string `json:"host"`
	} `json:"tag_service"`
	CalendarService struct {
		Port string `json:"port"`
		Host string `json:"host"`
	} `json:"calendar_service"`
}

func GetIP() string {
	name, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "error"
	}
	var realIp string

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				realIp = ipnet.IP.String()
			}
		}
	}
	if realIp != "" {
		return name + " " + realIp
	}

	return name + " unkownIP"
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("["+ GetIP()+ "] ")

	config := &Config{
		Database: struct {
			Host     string `json:"host"`
			Username string `json:"username"`
			Password string `json:"password"`
		}{Host: "couchbase://localhost", Username: "Administrator", Password: "Administrator"},
		Cookie64:    "f56931b180a522c21c17ba9274d07b0934a366bef09a7b43db9006a06b97bc3cf4a4e2c2ae77d11d3c38c6f17310eac0fbbb8da5eb1d1ae57eaef133abbf5ec8",
		Cookie32:    "20c2f69c5d03cd82deb2203bc9b1701fad3b7c3a7d3d045c43bb6994ccb44afd",
		ApiGateway:  struct{ Port string `json:"port"`; Host string `json:"host"` }{Port: "4430", Host: "apigateway:4430"},
		ApiService:  struct{ Port string `json:"port"`; Host string `json:"host"` }{Port: "50001", Host: "apiservice:50001"},
		HtmlService:  struct{ Port string `json:"port"`; Host string `json:"host"` }{Port: "50002", Host: "htmlservice:50002"},
		UserService: struct{ Port string `json:"port"`; Host string `json:"host"`  }{Port: "50003", Host: "userservice:50003"},
		MealService: struct{ Port string `json:"port"`; Host string `json:"host"`  }{Port: "50004", Host: "mealservice:50004"},
		TagService: struct{ Port string `json:"port"` ; Host string `json:"host"`  }{Port: "50005", Host: "tagservice:50005"},
		CalendarService: struct{ Port string `json:"port"` ; Host string `json:"host"`  }{Port: "50006", Host: "calendarservice:50006"},
	}

	configJson, _ := json.Marshal(config)

	log.Println("Config.json is up and running...")
	log.Fatal(http.ListenAndServe(":50000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		w.Write(configJson)
	})))
}
