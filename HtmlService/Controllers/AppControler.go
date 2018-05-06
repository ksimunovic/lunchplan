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
)

type Controller struct {
	ControllerName string
}

func getTemplate(method string, controllerName string) *template.Template {
	lp := filepath.Join("HtmlService", "templates", "layout.html")
	fp := filepath.Join("HtmlService", "templates", controllerName+"_"+method+".html")

	tmpl := template.New("home")

	funcMap := template.FuncMap{}
	funcMap["dict"] = dict
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
		fmt.Printf("\nRender Error: %v\n", err)
		return
	}
	w.Write(buf.Bytes())
}
