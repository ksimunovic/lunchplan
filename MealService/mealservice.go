package main

import (
	"encoding/json"
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
	ServedBy  string        `json:"served_by,omitempty" bson:"-"`
}
type Session struct {
	Id      bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Profile Profile       `json:"profile,omitempty"`
	Created time.Time     `json:"created,omitempty" bson:"created"`
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
	Id       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string        `json:"name,omitempty"`
	Profile  Profile       `json:"profile,omitempty"`
	ServedBy string        `json:"served_by,omitempty" bson:"-"`
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

func GetIP() string {
	name, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
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

type Server struct{}

func (s *Server) Negate(i int64, reply *int64) error {
	*reply = -i
	return nil
}

func (s *Server) Create(jsonData []byte, jsonResponse *[]byte) error {
	var data map[string]interface{}
	_ = json.Unmarshal(jsonData, &data)

	rpcData := map[string]interface{}{
		"sid": data["sid"].(string),
	}
	var profile Profile
	rpcResult := ServiceCallData("GetAccount", rpcData, LoadConfiguration().UserService.Host)
	temp0, _ := json.Marshal(rpcResult)
	_ = json.Unmarshal(temp0, &profile)

	var tags []Tag
	if data["tags"] != nil {
		for i := 0; i < len(data["tags"].([]interface{})); i++ {
			var tag Tag
			tagData := data["tags"].([]interface{})[i].(map[string]interface{})
			tagData["sid"] = data["sid"].(string)
			if tagData["id"] != nil && tagData["id"].(string) != "" {
				tagData["get_id"] = tagData["id"]
				tagResult := ServiceCallData("Read", tagData, LoadConfiguration().TagService.Host)
				temp1, _ := json.Marshal(tagResult)
				_ = json.Unmarshal(temp1, &tag)
			} else {
				tagResult := ServiceCallData("Create", tagData, LoadConfiguration().TagService.Host)
				temp1, _ := json.Marshal(tagResult)
				_ = json.Unmarshal(temp1, &tag)
			}
			tags = append(tags, tag)
		}
	}

	meal := Meal{
		Id:          bson.NewObjectId(),
		Title:       data["title"].(string),
		Description: data["description"].(string),
		Profile:     profile,
		Timestamp:   int(time.Now().Unix()),
		Tags:        tags,
	}

	sessionCopy := dbSession.Copy()
	defer sessionCopy.Close()

	c := sessionCopy.DB("MealService").C("meal")
	if err := c.Insert(meal); err != nil {
		panic(err)
	}

	//TODO update all tags async

	*jsonResponse, _ = json.Marshal(meal)
	return nil
}

func (s *Server) Read(jsonData []byte, jsonResponse *[]byte) error {
	var data map[string]interface{}
	_ = json.Unmarshal(jsonData, &data)

	rpcData := map[string]interface{}{
		"sid": data["sid"].(string),
	}
	var profile Profile
	rpcResult := ServiceCallData("GetAccount", rpcData, LoadConfiguration().UserService.Host)
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
		result.ServedBy = GetIP()
		*jsonResponse, _ = json.Marshal(result)
	}

	return nil
}

func (s *Server) Update(jsonData []byte, jsonResponse *[]byte) error {
	var data map[string]interface{}
	_ = json.Unmarshal(jsonData, &data)

	rpcData := map[string]interface{}{
		"sid": data["sid"].(string),
	}
	var profile Profile
	rpcResult := ServiceCallData("GetAccount", rpcData, LoadConfiguration().UserService.Host)
	temp0, _ := json.Marshal(rpcResult)
	_ = json.Unmarshal(temp0, &profile)

	if len(data["get_id"].(string)) != 12 && len(data["get_id"].(string)) != 24 {
		return errors.New("Invalid meal id in GET parameter")
	}

	sessionCopy := dbSession.Copy()
	defer sessionCopy.Close()

	c := sessionCopy.DB("MealService").C("meal")
	var meal Meal
	if err := c.Find(bson.M{"_id": bson.ObjectIdHex(data["get_id"].(string)), "profile._id": bson.ObjectIdHex(profile.Id.Hex())}).One(&meal); err != nil {
		return errors.New("Meal with id: " + data["get_id"].(string) + " from user: " + profile.Id.String() + " doesn't exist")
	}

	var tags []Tag
	if data["tags"] != nil {
		for i := 0; i < len(data["tags"].([]interface{})); i++ {
			var tag Tag
			tagData := data["tags"].([]interface{})[i].(map[string]interface{})
			tagData["sid"] = data["sid"].(string)
			if tagData["id"] != nil && tagData["id"].(string) != "" {
				tagData["get_id"] = tagData["id"]
				tagResult := ServiceCallData("Read", tagData, LoadConfiguration().TagService.Host)
				temp1, _ := json.Marshal(tagResult)
				_ = json.Unmarshal(temp1, &tag)
			} else {
				tagResult := ServiceCallData("Create", tagData, LoadConfiguration().TagService.Host)
				temp1, _ := json.Marshal(tagResult)
				_ = json.Unmarshal(temp1, &tag)
			}
			tags = append(tags, tag)
		}
	}

	updatedMeal := meal
	updatedMeal.Title = data["title"].(string)
	updatedMeal.Description = data["description"].(string)
	updatedMeal.Tags = tags
	/*
	data["id"] = data["get_id"]

	cleanUp := make(map[string]string)
	cleanUp["title"] = data["title"].(string)
	cleanUp["description"] = data["description"].(string)

	finalbody, err := json.Marshal(cleanUp)
	if err != nil {
		return err
	}
	var finalbodymap map[string]interface{}
	if err = json.Unmarshal(finalbody, &finalbodymap); err != nil {
		return err
	}
*/
	change := bson.M{"$set": updatedMeal}
	err := c.Update(meal, change)
	if err != nil {
		return err
	} else {
		*jsonResponse, _ = json.Marshal(make(map[string]interface{}))
	}

	return nil
}

func (s *Server) Delete(jsonData []byte, jsonResponse *[]byte) error {
	var data map[string]interface{}
	_ = json.Unmarshal(jsonData, &data)

	rpcData := map[string]interface{}{
		"sid": data["sid"].(string),
	}

	var profile Profile
	rpcResult := ServiceCallData("GetAccount", rpcData, LoadConfiguration().UserService.Host)
	temp0, _ := json.Marshal(rpcResult)
	_ = json.Unmarshal(temp0, &profile)

	if len(data["get_id"].(string)) != 12 && len(data["get_id"].(string)) != 24 {
		return errors.New("Invalid meal id in GET parameter")
	}

	sessionCopy := dbSession.Copy()
	defer sessionCopy.Close()

	c := sessionCopy.DB("MealService").C("meal")

	/*
	var result Meal
	if err := c.Find(bson.M{"_id": bson.ObjectIdHex(data["get_id"].(string)), "profile._id": bson.ObjectIdHex(profile.Id.Hex())}).One(&result); err != nil {
		return errors.New("Meal with id: " + data["get_id"].(string) + " from user: " + profile.Id.String() + " doesn't exist")
	}
	//TODO: Loop through all tags, if occurances in database == 1 then delete
	*/

	if err := c.Remove(bson.M{"_id": bson.ObjectIdHex(data["get_id"].(string)), "profile._id": bson.ObjectIdHex(profile.Id.Hex())}); err != nil {
		return errors.New("Unable to remove meal with id: " + data["get_id"].(string) + " from user: " + profile.Id.String())
	} else {
		*jsonResponse, _ = json.Marshal(make(map[string]interface{}))
	}

	return nil
}

func (s *Server) GetAllUserMeals(jsonData []byte, jsonResponse *[]byte) error {
	var data map[string]interface{}
	_ = json.Unmarshal(jsonData, &data)

	rpcData := map[string]interface{}{
		"sid": data["sid"].(string),
	}
	var profile Profile
	rpcResult := ServiceCallData("GetAccount", rpcData, LoadConfiguration().UserService.Host)
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

	for i := 0; i < len(results); i++ {
		results[i].ServedBy = GetIP()
	}

	*jsonResponse, _ = json.Marshal(results)
	return nil
}

func (s *Server) Suggest(jsonData []byte, jsonResponse *[]byte) error {
	var data map[string]interface{}
	_ = json.Unmarshal(jsonData, &data)

	rpcData := map[string]interface{}{
		"sid": data["sid"].(string),
	}
	var profile Profile
	rpcResult := ServiceCallData("GetAccount", rpcData, LoadConfiguration().UserService.Host)
	temp0, _ := json.Marshal(rpcResult)
	_ = json.Unmarshal(temp0, &profile)

	sessionCopy := dbSession.Copy()
	defer sessionCopy.Close()

	c := sessionCopy.DB("MealService").C("meal")
	var results []Meal
	help := []bson.M{{"$match": bson.M{"profile._id": bson.ObjectIdHex(profile.Id.Hex())}}, {"$sample": bson.M{"size": 1}}}
	pipe := c.Pipe(help)
	err := pipe.All(&results)
	if err != nil {
		panic(err)
	}

	if len(results) != 0 {
		*jsonResponse, _ = json.Marshal(results[0])
	} else {
		*jsonResponse, _ = json.Marshal(make(map[string]interface{}))
	}

	return nil
}

func ServiceCallData(method string, data map[string]interface{}, serviceHost string) map[string]interface{} {

	c, err := rpc.Dial("tcp", serviceHost)
	if err != nil {
		return nil
	}

	var rpcData []byte
	var result map[string]interface{}
	jsonData, _ := json.Marshal(data)
	err = c.Call("Server."+method, jsonData, &rpcData)
	if err != nil {
		return nil
	} else {
		_ = json.Unmarshal(rpcData, &result)
		return result
	}
}

var dbSession *mgo.Session

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("[" + GetIP() + "] ")

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
	log.Println("Meal Service RPC server online!")
	ln, err := net.Listen("tcp", ":"+LoadConfiguration().MealService.Port)
	if err != nil {
		log.Fatalln(err)
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
