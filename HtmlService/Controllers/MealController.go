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
	ServedBy  string        `json:"served_by,omitempty" bson:"-"`
}

type Meal struct {
	Id          bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string        `json:"title,omitempty"`
	Description string        `json:"description,omitempty"`
	Profile     Profile       `json:"profile,omitempty"`
	Timestamp   int           `json:"timestamp,omitempty"`
	Tags        []Tag         `json:"tags,omitempty"`
	ServedBy    string        `json:"served_by,omitempty" bson:"-"`
}
type Tag struct {
	Id   bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name string        `json:"name,omitempty"`
	Profile     Profile       `json:"profile,omitempty"`
	ServedBy    string        `json:"served_by,omitempty" bson:"-"`
}

func (c *Controller) Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rpcData := map[string]interface{}{
			"sid": r.Header.Get("Sid"),
		}

		if rpcData["sid"] == "" {
			println("HALLO mealcontroller nema sid")
		}

		var allUserMeals []Meal
		rpcResult := ServiceCallData("GetAllUserMeals", rpcData, LoadConfiguration().MealService.Port);
		if err := json.Unmarshal(rpcResult, &allUserMeals); err != nil {
			println(err.Error())
			return
		}
		allUserMealsJson, _ := json.Marshal(allUserMeals)

		var allUserTags []Tag
		rpcResult = ServiceCallData("GetAllUserTags", rpcData, LoadConfiguration().TagService.Port);
		if err := json.Unmarshal(rpcResult, &allUserTags); err != nil {
			println(err.Error())
			return
		}
		allUserTagsJson, _ := json.Marshal(allUserTags)

		vars := map[string]interface{}{
			"mealsJson": string(allUserMealsJson),
			"tagsJson": string(allUserTagsJson),
		}

		render(w, r, getTemplate(currentFunctionName(), c.ControllerName), "home", vars)
	}
}
