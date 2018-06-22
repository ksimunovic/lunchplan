package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"net/rpc"
	"time"
	"net/url"
	"github.com/gorilla/handlers"
	"regexp"
	"net"
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
}

var config Config

func LoadConfiguration() Config {
	if (Config{}) != config {
		return config
	}
	response, err := http.Get("http://configservice:50000")
	if err != nil {
		log.Fatalf("%s", err)
		log.Printf("Trying again in 5 seconds...")
		time.Sleep(5 * time.Second)
		return LoadConfiguration()
	} else {
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatalf("%s", err)
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

func ValidateApi(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {

				data := map[string]interface{}{
					"bearerToken": bearerToken[1],
				}

				c, err := rpc.Dial("tcp", LoadConfiguration().UserService.Host)
				if err != nil {
					log.Fatalln(err)
					return
				}

				var result []byte
				var profile map[string]interface{}
				jsonData, _ := json.Marshal(data)
				err = c.Call("Server.Validate", jsonData, &result)
				if err != nil {
					profile = make(map[string]interface{})
					profile["error"] = err.Error()
					json, _ := json.Marshal(profile)
					w.Write(json)
				} else {
					_ = json.Unmarshal(result, &profile)
					req.Header.Set("sid", profile["id"].(string))
					next(w, req)
				}
			}
		} else {
			req.Header.Set("sid", "")
			next(w, req)
		}
	})
}

func ValidateHtml(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cookie, _ := req.Cookie("sid")
		if cookie != nil && cookie.Value != "" {

			if req.RequestURI == "/logout" {
				c := &http.Cookie{
					Name:    "sid",
					Value:   "",
					Path:    "/",
					Expires: time.Unix(0, 0),

					HttpOnly: true,
				}

				http.SetCookie(w, c)
				http.Redirect(w, req, "/login", http.StatusSeeOther)
			}

			data := map[string]interface{}{
				"bearerToken": cookie.Value,
			}

			c, err := rpc.Dial("tcp", LoadConfiguration().UserService.Host)
			if err != nil {
				log.Fatalln(err)
				return
			}

			var result []byte
			var profile map[string]interface{}
			jsonData, _ := json.Marshal(data)
			err = c.Call("Server.Validate", jsonData, &result)
			if err != nil {
				c := &http.Cookie{
					Name:    "sid",
					Value:   "",
					Path:    "/",
					Expires: time.Unix(0, 0),

					HttpOnly: true,
				}

				http.SetCookie(w, c)
				http.Redirect(w, req, "/login", http.StatusSeeOther)
			} else {

				_ = json.Unmarshal(result, &profile)
				req.Header.Set("sid", profile["id"].(string))

				if req.RequestURI == "/login" {
					http.Redirect(w, req, "/", http.StatusSeeOther)
				} else {
					cookie := http.Cookie{Name: "sid", Path: "/", Value: profile["id"].(string), Expires: time.Now().Add(30 * time.Minute)}
					http.SetCookie(w, &cookie)
					next(w, req)
				}
			}

		} else {
			if req.RequestURI != "/login" && !strings.HasPrefix(req.RequestURI, "/static") {
				c := &http.Cookie{
					Name:    "sid",
					Value:   "",
					Path:    "/",
					Expires: time.Unix(0, 0),

					HttpOnly: true,
				}

				http.SetCookie(w, c)
				http.Redirect(w, req, "/login", http.StatusSeeOther)
			}
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
	/*
		Run Mac OS to get https://localhost/ to work
		echo "
		rdr pass inet proto tcp from any to any port 443 -> 127.0.0.1 port 4430
		" | sudo pfctl -ef -
	 */
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("["+ GetIP()+ "] ")

	u, _ := url.Parse("http://" + LoadConfiguration().ApiService.Host)
	apiProxy := httputil.NewSingleHostReverseProxy(u)

	h, _ := url.Parse("http://" + LoadConfiguration().HtmlService.Host)
	htmlProxy := httputil.NewSingleHostReverseProxy(h)

	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/api/{rest:.*}", ValidateApi(handler(apiProxy)))
	router.HandleFunc("/{rest:.*}", ValidateHtml(handler(htmlProxy)))

	go http.ListenAndServe(":80", http.HandlerFunc(redirect))

	log.Println("API Gateway service is up and running...")
	log.Fatal(http.ListenAndServeTLS(":4430", "certs/localhost.crt", "certs/localhost.key", logRequest(handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router))))
}

func redirect(w http.ResponseWriter, req *http.Request) {
	// remove/add not default ports from req.Host
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	log.Printf("redirect to: %s", target)
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Fatalln(err)
		}
		url := r.RequestURI
		re := regexp.MustCompile(`\r?\n`)
		log.Println("[REQUEST] ", url, " ", re.ReplaceAllString(string(requestDump), "; "))
		handler.ServeHTTP(w, r)
	})
}

