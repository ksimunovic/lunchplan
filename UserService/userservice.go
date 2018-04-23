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

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"database"`
	UserService struct {
		Port string `json:"port"`
	} `json:"user_service"`
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

type Server struct{}

func (s *Server) Negate(i int64, reply *int64) error {
	*reply = -i
	return nil
}

func (s *Server) Login(data map[string]interface{}, jsonResponse *[]byte) error {

	sessionCopy := dbSession.Copy()
	defer sessionCopy.Close()
	c := sessionCopy.DB("UserService").C("account")

	var account Account
	err := c.Find(bson.M{"email": data["email"].(string)}).One(&account)

	//TODO: Handle unknown email

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(data["password"].(string)))
	if err != nil {
		return err
	}

	i := bson.NewObjectId()
	session := Session{
		Id:        i,
		Profile:   account.Profile,
		CreatedAt: time.Now(),
	}

	c = sessionCopy.DB("UserService").C("session")
	if err := c.Insert(session); err != nil {
		panic(err)
	}

	var result map[string]interface{}
	result = make(map[string]interface{})
	result["sid"] = i.Hex()

	*jsonResponse, _ = json.Marshal(result)
	return nil
}

func (s *Server) Register(data map[string]interface{}, jsonResponse *[]byte) error {

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(data["password"].(string)), 10)
	i := bson.NewObjectId()
	profile := Profile{
		Id:        i,
		Firstname: data["firstname"].(string),
		Lastname:  data["lastname"].(string),
	}
	account := Account{
		Profile:  profile,
		Email:    data["email"].(string),
		Password: string(passwordHash),
	}

	//TODO: Missing email validation - ili?

	sessionCopy := dbSession.Copy()
	defer sessionCopy.Close()

	c := sessionCopy.DB("UserService").C("profile")
	if err := c.Insert(profile); err != nil {
		panic(err)
	}

	c = sessionCopy.DB("UserService").C("account")
	if err := c.Insert(account); err != nil {
		panic(err)
	}

	*jsonResponse, _ = json.Marshal(account)
	return nil
}

func (s *Server) GetAccount(data map[string]interface{}, jsonResponse *[]byte) error {
	sessionId := data["sid"].(string)
	sessionCopy := dbSession.Copy()
	defer sessionCopy.Close()
	c := sessionCopy.DB("UserService").C("session")

	var session Session
	_ = c.Find(bson.M{"_id": bson.ObjectIdHex(sessionId)}).One(&session)

	session.Profile.ServedBy = GetIP()

	*jsonResponse, _ = json.Marshal(session.Profile)
	return nil
}

func (s *Server) Blog(data map[string]interface{}, jsonResponse *[]byte) error {
	sessionId := data["sid"].(string)
	sessionCopy := dbSession.Copy()
	defer sessionCopy.Close()
	c := sessionCopy.DB("UserService").C("session")

	var session Session
	_ = c.Find(bson.M{"_id": bson.ObjectIdHex(sessionId)}).One(&session)

	temp0, _ := json.Marshal(data)

	var blog Blog
	_ = json.Unmarshal(temp0, &blog)
	blog.Type = "blog"
	blog.Profile = session.Profile
	blog.Timestamp = int(time.Now().Unix())

	c = sessionCopy.DB("UserService").C("blog")
	if err := c.Insert(blog); err != nil {
		panic(err)
	}

	*jsonResponse, _ = json.Marshal(blog)
	return nil
}

func (s *Server) Blogs(data map[string]interface{}, jsonResponse *[]byte) error {
	sessionId := data["sid"].(string)
	sessionCopy := dbSession.Copy()
	defer sessionCopy.Close()
	c := sessionCopy.DB("UserService").C("session")

	var session Session
	_ = c.Find(bson.M{"_id": bson.ObjectIdHex(sessionId)}).One(&session)

	c = sessionCopy.DB("UserService").C("blog")
	var results []Blog
	_ = c.Find(bson.M{"profile._id": bson.ObjectIdHex(session.Profile.Id.Hex())}).All(&results)

	*jsonResponse, _ = json.Marshal(results)
	return nil
}

func (s *Server) Validate(data map[string]interface{}, jsonResponse *[]byte) error {

	bearerToken := data["bearerToken"].(string)

	sessionCopy := dbSession.Copy()
	defer sessionCopy.Close()
	c := sessionCopy.DB("UserService").C("session")

	var session Session
	err := c.Find(bson.M{"_id": bson.ObjectIdHex(bearerToken)}).One(&session)
	if err != nil {
		return errors.New("Invalid session")
	}

	c.Update(bson.M{"_id": bson.ObjectIdHex(bearerToken)}, bson.M{"$set": bson.M{"createdat": time.Now()}})

	*jsonResponse, _ = json.Marshal(session)
	return nil
}

var dbSession *mgo.Session

func main() {

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{"localhost:27017"},
		Timeout:  60 * time.Second,
		Database: "UserService",
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
	fmt.Println("User Service RPC server online!")
	ln, err := net.Listen("tcp", ":"+LoadConfiguration().UserService.Port)
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
