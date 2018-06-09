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
	"gopkg.in/mgo.v2/bson"
	"errors"
	"gopkg.in/mgo.v2"
	"time"
)

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
	ServedBy    string        `json:"served_by,omitempty"`
	Tags        []Tag         `json:"tags,omitempty"`
}
type Tag struct {
	Id   bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name string        `json:"name,omitempty"`
	Profile     Profile       `json:"profile,omitempty"`
	ServedBy    string        `json:"served_by,omitempty" bson:"-"`
}

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"database"`
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
}

var config Config

func LoadConfiguration() Config {
	if (Config{}) != config {
		return config
	}
	response, err := http.Get("http://configservice:50000")
	if err != nil {
		fmt.Printf("%s; ", err)
		fmt.Println("Trying again in 5 seconds...")
		time.Sleep(5 * time.Second)
		return LoadConfiguration()
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

func (s *Server) Create(jsonData []byte, jsonResponse *[]byte) error {
	var data map[string]interface{}
	_ = json.Unmarshal(jsonData, &data)

	rpcData := map[string]interface{}{
		"sid": data["sid"].(string),
	}
	var profile Profile
	rpcResult := ServiceCallData("GetAccount", rpcData, LoadConfiguration().UserService.Host)
	if err := json.Unmarshal(rpcResult, &profile); err != nil {
		println(err.Error())
	}

	tag := Tag{
		Id:   bson.NewObjectId(),
		Name: data["name"].(string),
		Profile: profile,
	}

	sessionCopy := dbSession.Copy()
	defer sessionCopy.Close()

	c := sessionCopy.DB("TagService").C("tag")
	if err := c.Insert(tag); err != nil {
		panic(err)
	}

	*jsonResponse, _ = json.Marshal(tag)
	return nil
}

func (s *Server) Read(jsonData []byte, jsonResponse *[]byte) error {
	var data map[string]interface{}
	_ = json.Unmarshal(jsonData, &data)

	sessionCopy := dbSession.Copy()
	defer sessionCopy.Close()

	c := sessionCopy.DB("TagService").C("tag")
	var result Tag
	if err := c.Find(bson.M{"_id": bson.ObjectIdHex(data["get_id"].(string))}).One(&result); err != nil {
		return errors.New("Tag with id: " + data["get_id"].(string) + " doesn't exist")
	} else {
		*jsonResponse, _ = json.Marshal(result)
	}

	return nil
}

func (s *Server) Update(jsonData []byte, jsonResponse *[]byte) error {
	var data map[string]interface{}
	_ = json.Unmarshal(jsonData, &data)

	sessionCopy := dbSession.Copy()
	defer sessionCopy.Close()

	if len(data["get_id"].(string)) != 12 && len(data["get_id"].(string)) != 24 {
		return errors.New("Invalid tag id in GET parameter")
	}

	c := sessionCopy.DB("TagService").C("tag")
	var result Tag
	if err := c.Find(bson.M{"_id": bson.ObjectIdHex(data["get_id"].(string))}).One(&result); err != nil {
		return errors.New("Tag with id: " + data["get_id"].(string) + " doesn't exist")
	} else {
		*jsonResponse, _ = json.Marshal(make(map[string]interface{}))
	}

	if _, ok := data["get_id"]; !ok {
		return errors.New("Can't update without tag id")
	}

	finalbodymap := make(map[string]interface{})
	finalbodymap["name"] = data["name"]

	change := bson.M{"$set": finalbodymap}
	err := c.Update(result, change)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Delete(jsonData []byte, jsonResponse *[]byte) error {
	var data map[string]interface{}
	_ = json.Unmarshal(jsonData, &data)

	sessionCopy := dbSession.Copy()
	defer sessionCopy.Close()

	if len(data["get_id"].(string)) != 12 && len(data["get_id"].(string)) != 24 {
		return errors.New("Invalid tag id in GET parameter")
	}

	c := sessionCopy.DB("TagService").C("tag")
	if err := c.Remove(bson.M{"_id": bson.ObjectIdHex(data["get_id"].(string))}); err != nil {
		return errors.New("Unable to remove tag with id: " + data["get_id"].(string))
	} else {
		*jsonResponse, _ = json.Marshal(make(map[string]interface{}))
	}

	return nil
}

func (s *Server) GetAllUserTags(jsonData []byte, jsonResponse *[]byte) error {
	var data map[string]interface{}
	_ = json.Unmarshal(jsonData, &data)

	rpcData := map[string]interface{}{
		"sid": data["sid"].(string),
	}
	var profile Profile
	rpcResult := ServiceCallData("GetAccount", rpcData, LoadConfiguration().UserService.Host)
	if err := json.Unmarshal(rpcResult, &profile); err != nil {
		println(err.Error())
	}


	sessionCopy := dbSession.Copy()
	defer sessionCopy.Close()

	c := sessionCopy.DB("TagService").C("tag")
	var results []Tag
	if err := c.Find(bson.M{"profile._id": bson.ObjectIdHex(profile.Id.Hex())}).All(&results); err != nil {
		panic(err)
	}

	if len(results) == 0 {
		results = []Tag{}
	}

	for i := 0; i < len(results); i++ {
		results[i].ServedBy = GetIP()
	}

	*jsonResponse, _ = json.Marshal(results)
	return nil

	/*
	var meals []Meal
	rpcResult = ServiceCallData("GetAllUserMeals", rpcData, LoadConfiguration().MealService.Port);
	fmt.Println(string(rpcResult))
	if err := json.Unmarshal(rpcResult, &meals); err != nil {
		println(err.Error())
	}


	var results []Tag
	if profile.Id == "" {
		results = []Tag{}
	} else {
		for i := 0; i < len(meals); i++ {
			meal := meals[i]
			for t := 0; t < len(meal.Tags); t++ {
				results = append(results, meal.Tags[t])
			}
		}
	}*/
}

func ServiceCallData(method string, data map[string]interface{}, serviceHost string) []byte {

	c, err := rpc.Dial("tcp", serviceHost)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	if data["sid"] == "" {
		fmt.Println("Missing sid from rpc data request")
		return nil
	}

	var rpcData []byte
	jsonData, _ := json.Marshal(data)
	err = c.Call("Server."+method, jsonData, &rpcData)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	} else {
		return rpcData
	}
}

var dbSession *mgo.Session

func main() {

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{"mongodb:27017"},
		Timeout:  60 * time.Second,
		Database: "admin",
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
	fmt.Println("Tag Service RPC server online!")
	ln, err := net.Listen("tcp", ":"+LoadConfiguration().TagService.Port)
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
