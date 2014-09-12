package main

import (
	"io/ioutil"
	"log"
	"fmt"
	"encoding/json"
	//"net/http"
)

type RequestInfo struct {
	Url string
	Secure bool
	Interval int
	Script string
	Email bool
	Email_address string
	Log bool
	Logfile string
}

var reqinfo []RequestInfo

func json_to_requestinfo(j []byte) (r []RequestInfo, err error) {
	err = json.Unmarshal(j, &r)
	return r, err
}

func getRequestInfo(f string) (r []RequestInfo, err error) {
	var the_json []byte
	the_json, err = ioutil.ReadFile(f)
	if err == nil {
		log.Print("JSON file loaded. OK")
		r, err = json_to_requestinfo(the_json)
	}
	return r, err
}

func main() {
	// TODO: json file name via argument
	// TODO: list of json files, via argument
	json_file := "./files/gocontrol.json"
	// get the requests to monitor

	res, err := getRequestInfo(json_file)
	if  err != nil {
		log.Fatal("Fatal error loading the JSON data,", err)
	}
	// fmt.Printf("Results: %+v\n", res)

	// create go routines for each request

}

