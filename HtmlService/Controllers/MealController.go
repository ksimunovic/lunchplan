package Controller

import (
	"net/http"
)

var Meal = Controller{ControllerName: "meal"}

func (c *Controller) Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fullData := map[string]interface{}{
			"User": "Dewey",
		}
		render(w, r, getTemplate(currentFunctionName(), c.ControllerName), "home", fullData)
	}
}
