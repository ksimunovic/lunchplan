package Controller

import (
	"net/http"
)

var User = Controller{ControllerName: "user"}

func (c *Controller) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fullData := map[string]interface{}{
			"User": "Dewey",
		}
		render(w, r, getTemplate(currentFunctionName(), c.ControllerName), "home", fullData)
	}
}
