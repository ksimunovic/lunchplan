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
	"time"
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
	} `json:"api_service"`
	HtmlService struct {
		Port string `json:"port"`
	} `json:"html_service"`
	UserService struct {
		Port string `json:"port"`
	} `json:"user_service"`
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

func ValidateApi(next http.HandlerFunc) http.HandlerFunc {
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
					Name:     "sid",
					Value:    "",
					Path:     "/",
					Expires: time.Unix(0, 0),

					HttpOnly: true,
				}

				http.SetCookie(w, c)
				http.Redirect(w, req, "/login", http.StatusSeeOther)
			}

			data := map[string]interface{}{
				"bearerToken": cookie.Value,
			}

			c, err := rpc.Dial("tcp", "127.0.0.1:"+LoadConfiguration().UserService.Port)
			if err != nil {
				fmt.Println(err)
				return
			}

			var result []byte
			var profile map[string]interface{}
			jsonData, _ := json.Marshal(data)
			err = c.Call("Server.Validate", jsonData, &result)
			if err != nil {
				c := &http.Cookie{
					Name:     "sid",
					Value:    "",
					Path:     "/",
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
					Name:     "sid",
					Value:    "",
					Path:     "/",
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
	logPath := "development.log"

	openLogFile(logPath)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	u, _ := url.Parse("http://localhost:" + LoadConfiguration().ApiService.Port)
	apiProxy := httputil.NewSingleHostReverseProxy(u)

	h, _ := url.Parse("http://localhost:" + LoadConfiguration().HtmlService.Port)
	htmlProxy := httputil.NewSingleHostReverseProxy(h)

	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/api/{rest:.*}", ValidateApi(handler(apiProxy)))
	router.HandleFunc("/{rest:.*}", ValidateHtml(handler(htmlProxy)))

	fmt.Println("API Gateway is up and running...")

	go http.ListenAndServe(":80", http.HandlerFunc(redirect))
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
			log.Println(err)
		}
		url := r.RequestURI
		log.Println("[REQUEST] ", url, " ", string(requestDump))
handler.ServeHTTP(w, r)
		/*
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, r)

		dump, err := httputil.DumpResponse(rec.Result(), false)
		if err != nil {
			log.Fatal(err)
		}
		log.Println( "[RESPONSE] ", url, " ", string(dump))

		// we copy the captured response headers to our new response
		for k, v := range rec.Header() {
			w.Header()[k] = v
		}

		// grab the captured response body
		data := rec.Body.Bytes()

		w.Write(data)
		*/
	})
}

func openLogFile(logfile string) {
	if logfile != "" {
		err := os.Truncate(logfile, 100)
		if err != nil {
			log.Fatal(err)
		}

		lf, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE, 0640)

		if err != nil {
			log.Fatal("OpenLogfile: os.OpenFile:", err)
		}

		log.SetOutput(lf)
	}
}
