package Controller

import (
	"net/http"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
)

var MealController = Controller{ControllerName: "meal"}

type Profile struct {
	Id        bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Firstname string        `json:"firstname,omitempty"`
	Lastname  string        `json:"lastname,omitempty"`
	ServedBy  string        `json:"served_by,omitempty"`
}

type Meal struct {
	Id          bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string        `json:"title,omitempty"`
	Description string        `json:"description,omitempty"`
	Profile     Profile       `json:"profile,omitempty"`
	Timestamp   int           `json:"timestamp,omitempty"`
	ServedBy    string        `json:"served_by,omitempty"`
}

func (c *Controller) Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rpcData := map[string]interface{}{
			"sid": r.Header.Get("Sid"),
		}

		if rpcData["sid"] == "" {
			println("HALLO mealcontroller nema sid")
		}
/*
		rpcResult := ServiceCallData("GetAllUserMeals", rpcData, LoadConfiguration().MealService.Port);
		rawJson, _ := json.Marshal(rpcResult)
		vars := map[string]interface{}{
			"mealsJson": string(rawJson),
		}*/

		var result []Meal
		rpcResult := ServiceCallData("GetAllUserMeals", rpcData, LoadConfiguration().MealService.Port);
		if err := json.Unmarshal(rpcResult, &result); err != nil {
			println(err.Error())
			return
		}

		rawJson, _ := json.Marshal(result)

		vars := map[string]interface{}{
			"mealsJson": string(rawJson),
		}

		render(w, r, getTemplate(currentFunctionName(), c.ControllerName), "home", vars)
	}
}
