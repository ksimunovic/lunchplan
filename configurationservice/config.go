package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func LoadConfiguration(w http.ResponseWriter, r *http.Request) {
	config, _ := ioutil.ReadFile("configurationservice/config.json")
	b := bytes.NewBuffer(config)
	w.Header().Set("Content-type", "application/json")
	if _, err := b.WriteTo(w); err != nil {
		fmt.Fprintf(w, "%s", err)
	}
}

func main() {
	http.HandleFunc("/", LoadConfiguration)
	fmt.Println("Serving config.json...")
	if err := http.ListenAndServe(":50000", nil); err != nil {
		panic(err)
	}
}
