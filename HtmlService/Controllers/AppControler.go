package Controller

import (
	"net/http"
	"html/template"
	"errors"
	"path/filepath"
	"runtime"
	"strings"
	"regexp"
	"bytes"
	"fmt"
	"net/rpc"
	"encoding/json"
	"io/ioutil"
	"os"
	"log"
	"time"
)

type Controller struct {
	ControllerName string
}

type Config struct {
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

var Configuration Config

func LoadConfiguration() Config {
	if (Config{}) != Configuration {
		return Configuration
	}
	response, err := http.Get("http://configservice:50000/")
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
func getTemplate(method string, controllerName string) *template.Template {
	dir := ""
	if _, err := os.Stat("HtmlService"); err == nil {
		dir = "HtmlService"
	}

	lp := filepath.Join(dir, "templates", "layout.html")
	fp := filepath.Join(dir, "templates", controllerName+"_"+method+".html")

	tmpl := template.New("home")

	funcMap := template.FuncMap{}
	funcMap["dict"] = dict
	funcMap["translate"] = translate
	tmpl.Funcs(funcMap)

	tmpl, _ = tmpl.ParseFiles(lp, fp)
	return tmpl
}

// Thank you tux21b
func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func currentFunctionName() string {
	pc, _, _, _ := runtime.Caller(1)
	nameFull := runtime.FuncForPC(pc).Name()
	re := regexp.MustCompile(`\/(.+)\.(\w+)\.(\w+)$`)
	match := re.FindStringSubmatch(nameFull)
	return strings.ToLower(match[2])
}

// Render a template, or server error.
func render(w http.ResponseWriter, r *http.Request, tpl *template.Template, name string, data interface{}) {
	buf := new(bytes.Buffer)
	if err := tpl.ExecuteTemplate(buf, name, data); err != nil {
		log.Fatalf("\nRender Error: %v\n", err)
		return
	}
	w.Write(buf.Bytes())
}

func ServiceCallData(method string, data map[string]interface{}, serviceHost string) []byte {

	c, err := rpc.Dial("tcp", serviceHost)
	if err != nil {
		log.Fatalln(err.Error())
		return nil
	}

	if data["sid"] == "" {
		log.Fatalln("Missing sid from rpc data request")
		return nil
	}

	var rpcData []byte
	jsonData, _ := json.Marshal(data)
	err = c.Call("Server."+method, jsonData, &rpcData)

	if err != nil {
		log.Fatalln(err.Error())
		return nil
	} else {
		return rpcData
	}
}

func translate(key string, args ...interface{}) string {
	var locale = "en"
	if v, ok := Locales[locale]; ok {
		if v2, ok := v[key]; ok {
			if len(args) > 0 {
				return fmt.Sprintf(v2, args...)
			}
			return v2
		}
	}
	return ""
}
