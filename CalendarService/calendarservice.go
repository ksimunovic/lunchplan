package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"

	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"errors"
)

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
type Calendar struct {
	Id   bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Date time.Time     `json:"date,omitempty"`
	Meal Meal          `json:"meal,omitempty"`
}

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"database"`
	UserService struct {
		Port string `json:"port"`
	} `json:"user_service"`
	MealService struct {
		Port string `json:"port"`
	} `json:"meal_service"`
	CalendarService struct {
		Port string `json:"port"`
	} `json:"calendar_service"`
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

type Server struct{}

func (s *Server) Negate(i int64, reply *int64) error {
	*reply = -i
	return nil
}

func GetIP() string {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
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

func (s *Server) Create(data map[string]interface{}, jsonResponse *[]byte) error {
	rpcData := map[string]interface{}{
		"sid":    data["sid"].(string),
		"get_id": data["meal_id"].(string),
	}
	var meal Meal
	rpcResult := ServiceCallData("Read", rpcData, LoadConfiguration().MealService.Port)
	temp0, _ := json.Marshal(rpcResult)
	_ = json.Unmarshal(temp0, &meal)

	if meal.Id == "" {
		return errors.New("Meal with id " + meal.Id.String() + " doesn't exist")
	}

	date, err := time.Parse("2006-01-02", data["date"].(string))
	if err != nil {
		return errors.New("Date must be in format YYYY-MM-DD, given: " + data["date"].(string))
	}

	calendar := Calendar{
		Id:   bson.NewObjectId(),
		Date: date,
		Meal: meal,
	}

	sessionCopy := dbSession.Copy()
	defer sessionCopy.Close()

	c := sessionCopy.DB("CalendarService").C("calendar")
	if err := c.Insert(calendar); err != nil {
		panic(err)
	}

	*jsonResponse, _ = json.Marshal(calendar)
	return nil
}

func (s *Server) Read(data map[string]interface{}, jsonResponse *[]byte) error {

	sessionCopy := dbSession.Copy()
	defer sessionCopy.Close()

	c := sessionCopy.DB("CalendarService").C("calendar")
	var result Calendar
	if err := c.Find(bson.M{"_id": bson.ObjectIdHex(data["get_id"].(string))}).One(&result); err != nil {
		return errors.New("Calendar with id: " + data["get_id"].(string) + " doesn't exist")
	} else {
		*jsonResponse, _ = json.Marshal(result)
	}

	return nil
}

func (s *Server) Update(data map[string]interface{}, jsonResponse *[]byte) error {

	sessionCopy := dbSession.Copy()
	defer sessionCopy.Close()

	if len(data["get_id"].(string)) != 12 && len(data["get_id"].(string)) != 24 {
		return errors.New("Invalid calendar id in GET parameter")
	}

	c := sessionCopy.DB("CalendarService").C("calendar")
	var result Calendar
	if err := c.Find(bson.M{"_id": bson.ObjectIdHex(data["get_id"].(string))}).One(&result); err != nil {
		return errors.New("Calendar with id: " + data["get_id"].(string) + " doesn't exist")
	} else {
		*jsonResponse, _ = json.Marshal(make(map[string]interface{}))
	}

	if _, ok := data["get_id"]; !ok {
		return errors.New("Can't update without calendar id")
	}

	data["id"] = data["get_id"]

	date, err := time.Parse("2006-01-02", data["date"].(string))
	if err != nil {
		return errors.New("Date must be in format YYYY-MM-DD, given: " + data["date"].(string))
	}

	finalbodymap := make( map[string]interface{})
	finalbodymap["date"] = date

	change := bson.M{"$set": finalbodymap}
	err = c.Update(result, change)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Delete(data map[string]interface{}, jsonResponse *[]byte) error {

	sessionCopy := dbSession.Copy()
	defer sessionCopy.Close()

	if len(data["get_id"].(string)) != 12 && len(data["get_id"].(string)) != 24 {
		return errors.New("Invalid calendar in GET parameter")
	}

	c := sessionCopy.DB("CalendarService").C("calendar")
	if err := c.Remove(bson.M{"_id": bson.ObjectIdHex(data["get_id"].(string))}); err != nil {
		return errors.New("Unable to remove calendar with id: " + data["get_id"].(string))
	} else {
		*jsonResponse, _ = json.Marshal(make(map[string]interface{}))
	}

	return nil
}


func (s *Server) GetAllUserCalendars(data map[string]interface{}, jsonResponse *[]byte) error {

	rpcData := map[string]interface{}{
		"sid": data["sid"].(string),
	}
	var profile Profile
	rpcResult := ServiceCallData("GetAccount", rpcData, LoadConfiguration().UserService.Port)
	temp0, _ := json.Marshal(rpcResult)
	_ = json.Unmarshal(temp0, &profile)

	sessionCopy := dbSession.Copy()
	defer sessionCopy.Close()

	c := sessionCopy.DB("CalendarService").C("calendar")
	var results []Calendar
	if err := c.Find(bson.M{"meal.profile._id": bson.ObjectIdHex(profile.Id.Hex())}).All(&results); err != nil {
		panic(err)
	}

	if len(results) == 0 {
		results = []Calendar{}
	}

	*jsonResponse, _ = json.Marshal(results)
	return nil
}


func ServiceCallData(method string, data map[string]interface{}, servicePort string) map[string]interface{} {

	c, err := rpc.Dial("tcp", "127.0.0.1:"+servicePort)
	if err != nil {
		return nil
	}

	var rpcData []byte
	var result map[string]interface{}
	err = c.Call("Server."+method, data, &rpcData)
	if err != nil {
		return nil
	} else {
		_ = json.Unmarshal(rpcData, &result)
		return result
	}
}

var dbSession *mgo.Session

func main() {

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{"localhost:27017"},
		Timeout:  60 * time.Second,
		Database: "CalendarService",
		Username: "root",
		Password: "root",
	}

	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}

	defer mongoSession.Close()
	mongoSession.SetMode(mgo.Monotonic, true)

	dbSession = mongoSession.Copy()
	defer mongoSession.Close()

	rpc.Register(new(Server))
	fmt.Println("Calendar Service RPC server online!")
	ln, err := net.Listen("tcp", ":"+LoadConfiguration().CalendarService.Port)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(c)
	}
}
