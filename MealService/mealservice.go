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
)

type Account struct {
	Type     string `json:"type,omitempty"`
	Pid      string `json:"pid,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}
type Profile struct {
	Type      string `json:"type,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
}
type Session struct {
	Type string `json:"type,omitempty"`
	Pid  string `json:"pid,omitempty"`
}
type Blog struct {
	Type      string `json:"type,omitempty"`
	Pid       string `json:"pid,omitempty"`
	Title     string `json:"title,omitempty"`
	Content   string `json:"content,omitempty"`
	Timestamp int    `json:"timestamp,omitempty"`
}
type Meal struct {
	Type        string `json:"type,omitempty"`
	Pid         string `json:"pid,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"content,omitempty"`
	Timestamp   int    `json:"timestamp,omitempty"`
}

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"database"`
	UserService struct {
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
	var account Account
	_, err := bucket.Get(data["email"].(string), &account)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(data["password"].(string)))
	if err != nil {
		return err
	}
	session := Session{
		Type: "session",
		Pid:  account.Pid,
	}
	var result map[string]interface{}
	result = make(map[string]interface{})

	temp1, _ := uuid.NewV4()
	result["pid"] = temp1.String()
	_, err = bucket.Insert(result["pid"].(string), &session, 3600)
	if err != nil {
		return err
	}

	*jsonResponse, _ = json.Marshal(result)
	return nil
}

func (s *Server) Create(data map[string]interface{}, jsonResponse *[]byte) error {

	temp1, _ := uuid.NewV4()
	id := temp1.String()
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(data["password"].(string)), 10)
	account := Account{
		Type:     "account",
		Pid:      id,
		Email:    data["email"].(string),
		Password: string(passwordHash),
	}
	profile := Profile{
		Type:      "profile",
		Firstname: data["firstname"].(string),
		Lastname:  data["lastname"].(string),
	}

	_, err := bucket.Insert(id, profile, 0)
	if err != nil {
		return err
	}
	_, err = bucket.Insert(data["email"].(string), account, 0)
	if err != nil {
		return err
	}

	*jsonResponse, _ = json.Marshal(account)
	return nil
}

func (s *Server) GetAccount(data map[string]interface{}, jsonResponse *[]byte) error {
	pid := data["pid"].(string)
	var profile Profile
	_, err := bucket.Get(pid, &profile)
	if err != nil {
		return err
	}

	//TODO: Add served by for load balancing testing

	*jsonResponse, _ = json.Marshal(profile)
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

	var session Session
	_, err := bucket.Get(bearerToken, &session)
	if err != nil {
		return err
	}
	bucket.Touch(bearerToken, 0, 3600)

	*jsonResponse, _ = json.Marshal(session)
	return nil
}

var bucket *gocb.Bucket

func main() {
	cluster, _ := gocb.Connect(LoadConfiguration().Database.Host)
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: LoadConfiguration().Database.Username,
		Password: LoadConfiguration().Database.Password,
	})
	bucket, _ = cluster.OpenBucket("MealService", "")
	rpc.Register(new(Server))
	fmt.Println("Meal Service RPC server online!")
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
