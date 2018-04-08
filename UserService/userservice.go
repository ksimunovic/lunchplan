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

	"github.com/couchbase/gocb"
	"github.com/satori/go.uuid"
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
}
type Session struct {
	Id        bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Profile   Profile       `json:"profile,omitempty"`
	CreatedAt time.Time     `json:”created_at,omitempty” bson:"createdAt"`
}
type Blog struct {
	Type      string `json:"type,omitempty"`
	Pid       string `json:"pid,omitempty"`
	Title     string `json:"title,omitempty"`
	Content   string `json:"content,omitempty"`
	Timestamp int    `json:"timestamp,omitempty"`
}

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"database"`
	MealService struct {
		Port string `json:"port"`
	} `json:"UserService"`
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

	//TODO: Add served_by field for load balancing testing

	*jsonResponse, _ = json.Marshal(session.Profile)
	return nil
}

func (s *Server) Blog(data map[string]interface{}, jsonResponse *[]byte) error {

	temp0, _ := json.Marshal(data)

	var blog Blog
	_ = json.Unmarshal(temp0, &blog)
	blog.Type = "blog"

	blog.Timestamp = int(time.Now().Unix())
	temp1, _ := uuid.NewV4()
	_, err := bucket.Insert(temp1.String(), blog, 0)
	if err != nil {
		return err
	}

	*jsonResponse, _ = json.Marshal(blog)
	return nil
}

func (s *Server) Blogs(data map[string]interface{}, jsonResponse *[]byte) error {
	pid := data["pid"].(string)

	var n1qlParams []interface{}
	n1qlParams = append(n1qlParams, pid)
	query := gocb.NewN1qlQuery("SELECT `" + bucket.Name() + "`.* FROM `" + bucket.Name() + "` WHERE type = 'blog' AND pid = $1")
	query.Consistency(gocb.RequestPlus)
	rows, err := bucket.ExecuteN1qlQuery(query, n1qlParams)
	if err != nil {
		return err
	}

	var row map[string]interface{}
	var result []map[string]interface{}
	for rows.Next(&row) {
		result = append(result, row)
		row = make(map[string]interface{})
	}
	rows.Close()

	if result == nil {
		result = make([]map[string]interface{}, 0)
	}

	*jsonResponse, _ = json.Marshal(result)
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

var bucket *gocb.Bucket
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
