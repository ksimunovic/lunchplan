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

type Account struct {
	Profile  Profile `json:"profile,omitempty"`
	Email    string  `json:"email,omitempty"`
	Password string  `json:"password,omitempty"`
}
type Profile struct {
	Id        bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Firstname string        `json:"firstname,omitempty"`
	Lastname  string        `json:"lastname,omitempty"`
	ServedBy  string        `json:"served_by,omitempty"`
}
type Session struct {
	Id        bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Profile   Profile       `json:"profile,omitempty"`
	CreatedAt time.Time     `json:"created_at,omitempty" bson:"createdAt"`
}
type Blog struct {
	Type      string  `json:"type,omitempty"`
	Profile   Profile `json:"profile,omitempty"`
	Title     string  `json:"title,omitempty"`
	Content   string  `json:"content,omitempty"`
	Timestamp int     `json:"timestamp,omitempty"`
}

type Meal struct {
	Id          bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string        `json:"title,omitempty"`
	Description string        `json:"description,omitempty"`
	Profile     Profile       `json:"profile,omitempty"`
	Timestamp   int           `json:"timestamp,omitempty"`
	ServedBy    string        `json:"served_by,omitempty"`
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
		"sid": data["sid"].(string),
	}
	var profile Profile
	rpcResult := UserService("GetAccount", rpcData)
	temp0, _ := json.Marshal(rpcResult)
	_ = json.Unmarshal(temp0, &profile)

	meal := Meal{
		Title:       "Saft i tijesto",
		Description: "Svinjsko mljeveno meso i tjestenina",
		Profile:     profile,
		Timestamp:   int(time.Now().Unix()),
		ServedBy:    GetIP(),
	}

	sessionCopy := dbSession.Copy()
	defer sessionCopy.Close()

	c := sessionCopy.DB("MealService").C("meal")
	if err := c.Insert(meal); err != nil {
		panic(err)
	}

	*jsonResponse, _ = json.Marshal(meal)
	return nil
}

func (s *Server) Read(data map[string]interface{}, jsonResponse *[]byte) error {

	rpcData := map[string]interface{}{
		"sid": data["sid"].(string),
	}
	var profile Profile
	rpcResult := UserService("GetAccount", rpcData)
	temp0, _ := json.Marshal(rpcResult)
	_ = json.Unmarshal(temp0, &profile)

	sessionCopy := dbSession.Copy()
	defer sessionCopy.Close()

	c := sessionCopy.DB("MealService").C("meal")
	var result Meal
	if err := c.Find(bson.M{"_id": bson.ObjectIdHex(data["get_id"].(string)), "profile._id": bson.ObjectIdHex(profile.Id.Hex())}).One(&result); err != nil {
		return errors.New("Meal with id: " + data["get_id"].(string) + " from user: " + profile.Id.String() + " doesn't exist")
	} else {
		*jsonResponse, _ = json.Marshal(result)
	}

	return nil
}

func (s *Server) Update(data map[string]interface{}, jsonResponse *[]byte) error {

	rpcData := map[string]interface{}{
		"sid": data["sid"].(string),
	}
	var profile Profile
	rpcResult := UserService("GetAccount", rpcData)
	temp0, _ := json.Marshal(rpcResult)
	_ = json.Unmarshal(temp0, &profile)

	sessionCopy := dbSession.Copy()
	defer sessionCopy.Close()

	if len(data["get_id"].(string)) != 12 && len(data["get_id"].(string)) != 24 {
		return errors.New("Invalid meal id in GET parameter")
	}

	c := sessionCopy.DB("MealService").C("meal")
	var result Meal
	if err := c.Find(bson.M{"_id": bson.ObjectIdHex(data["get_id"].(string)), "profile._id": bson.ObjectIdHex(profile.Id.Hex())}).One(&result); err != nil {
		return errors.New("Meal with id: " + data["get_id"].(string) + " from user: " + profile.Id.String() + " doesn't exist")
	} else {
		*jsonResponse, _ = json.Marshal(make(map[string]interface{}))
	}

	if len(data) < 1{
		return errors.New("Empty data to update")
	} else if _, ok := data["get_id"]; !ok {
		return errors.New("Can't update without meal id")
	}

	data["id"] = data["get_id"]

	finalbody, err := json.Marshal(data)
	if err != nil {
		return err
	}
	var finalbodymap map[string]interface{}
	if err = json.Unmarshal(finalbody, &finalbodymap); err != nil{
		return err
	}

	change := bson.M{"$set": finalbodymap}
	err = c.Update(result, change)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Delete(data map[string]interface{}, jsonResponse *[]byte) error {

	rpcData := map[string]interface{}{
		"sid": data["sid"].(string),
	}
	var profile Profile
	rpcResult := UserService("GetAccount", rpcData)
	temp0, _ := json.Marshal(rpcResult)
	_ = json.Unmarshal(temp0, &profile)

	if len(data["get_id"].(string)) != 12 && len(data["get_id"].(string)) != 24 {
		return errors.New("Invalid meal id in GET parameter")
	}

	sessionCopy := dbSession.Copy()
	defer sessionCopy.Close()

	c := sessionCopy.DB("MealService").C("meal")
	if err := c.Remove(bson.M{"_id": bson.ObjectIdHex(data["get_id"].(string)), "profile._id": bson.ObjectIdHex(profile.Id.Hex())}); err != nil {
		return errors.New("Unable to remove meal with id: " + data["get_id"].(string) + " from user: " + profile.Id.String())
	} else {
		*jsonResponse, _ = json.Marshal(make(map[string]interface{}))
	}

	return nil
}

func (s *Server) GetAllMeals(data map[string]interface{}, jsonResponse *[]byte) error {

	rpcData := map[string]interface{}{
		"sid": data["sid"].(string),
	}
	var profile Profile
	rpcResult := UserService("GetAccount", rpcData)
	temp0, _ := json.Marshal(rpcResult)
	_ = json.Unmarshal(temp0, &profile)

	sessionCopy := dbSession.Copy()
	defer sessionCopy.Close()

	c := sessionCopy.DB("MealService").C("meal")
	var results []Meal
	if err := c.Find(bson.M{"profile._id": bson.ObjectIdHex(profile.Id.Hex())}).All(&results); err != nil {
		panic(err)
	}

	if len(results) == 0 {
		results = []Meal{}
	}

	*jsonResponse, _ = json.Marshal(results)
	return nil
}

func UserService(method string, data map[string]interface{}) map[string]interface{} {

	c, err := rpc.Dial("tcp", "127.0.0.1:"+LoadConfiguration().UserService.Port)
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
		Database: "MealService",
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
	fmt.Println("Meal Service RPC server online!")
	ln, err := net.Listen("tcp", ":"+LoadConfiguration().MealService.Port)
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
