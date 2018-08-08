package Controller

import (
	"net/http"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"time"
	"log"
)

type Calendar struct {
	Id       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Date     time.Time     `json:"date,omitempty"`
	Meal     Meal          `json:"meal,omitempty"`
	Start    string        `json:"start,omitempty" bson:"-"`
	End      string        `json:"end,omitempty" bson:"-"`
	Title      string        `json:"title,omitempty" bson:"-"`
	ServedBy string        `json:"served_by,omitempty" bson:"-"`
}

var CalendarController = Controller{ControllerName: "calendar"}

func (c *Controller) ShowCalendar() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rpcData := map[string]interface{}{
			"sid": r.Header.Get("Sid"),
		}

		if rpcData["sid"] == "" {
			log.Println("HALLO calendarcontroller nema sid")
		}

		var allUserCalendars []Calendar
		rpcResult := ServiceCallData("GetAllUserCalendars", rpcData, LoadConfiguration().CalendarService.Host);
		if err := json.Unmarshal(rpcResult, &allUserCalendars); err != nil {
			log.Println(err.Error())
			return
		}
		allUserCalendarsJson, _ := json.Marshal(allUserCalendars)

		var allUserMeals []Meal
		rpcResult = ServiceCallData("GetAllUserMeals", rpcData, LoadConfiguration().MealService.Host);
		if err := json.Unmarshal(rpcResult, &allUserMeals); err != nil {
			log.Println(err.Error())
			return
		}
		allUserMealsJson, _ := json.Marshal(allUserMeals)

		vars := map[string]interface{}{
			"calendarsJson": string(allUserCalendarsJson),
			"mealsJson": string(allUserMealsJson),
		}

		render(w, r, getTemplate(currentFunctionName(), c.ControllerName), "home", vars)
	}
}
