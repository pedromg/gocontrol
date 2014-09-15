package main

import (
	"io/ioutil"
	"log"
	"fmt"
	"sync"
	"time"
	"bytes"
	"os/exec"
	"encoding/json"
	"net/http"
)

type RequestInfo struct {
	Name string
	Url string
	Secure bool
	Interval int
	StatusCode []int
	MaxAlerts int
	Script string
	DelayedBy int
	Email bool
	EmailAddress string
	Log bool
	Logfile string
}

var reqinfo []RequestInfo

func json_to_requestinfo(j []byte) (r []RequestInfo, err error) {
	err = json.Unmarshal(j, &r)
	return r, err
}

func GetRequestInfo(f string) (r []RequestInfo, err error) {
	var the_json []byte
	the_json, err = ioutil.ReadFile(f)
	if err == nil {
		log.Print("JSON file loaded. OK")
		r, err = json_to_requestinfo(the_json)
	}
	return r, err
}

func sendEmail(elem RequestInfo, err error) {
	log.Print("EMAIL ALERT FOR ", elem.Name, " !! ", err )
}

func sendLog(elem RequestInfo, err error) {
	log.Print("LOG ALERT FOR ", elem.Name, " !! ", err )
}

func execScript(elem RequestInfo) (err error) {
	log.Print("EXEC SCRIPT FOR ", elem.Name)
	cmd := exec.Command("sh", elem.Script)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	return err
}

func startWorker(elem RequestInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	var url string
	alert		:= false
	delayCounter	:= 0
	maxAlertCounter	:= 0
	if elem.Secure {
		url = "https://" + elem.Url
	} else {
		url = "http://" + elem.Url
	}
	log.Print("Started monitoring >> ", elem.Url, " at a ", elem.Interval, "s interval")

	for {
		// request is sent
		r, err := http.Head(url)

		// check for errors and response. If no errors, check if response status code is in the codes slice.
		if err != nil {
			alert = true
		} else {
			for _, sc := range elem.StatusCode {
				if sc == r.StatusCode {
					alert = true
				}
			}
		}

		// delay script execution counter.
		// alerts will be sent if the maxalert is not reached.
		// maxAlertCounter is incremented or reset. 
		if alert {
			delayCounter += 1
			maxAlertCounter += 1
		} else {
			maxAlertCounter = 0
			delayCounter = 0
		}
		// script execution
		if alert && delayCounter > elem.DelayedBy {
			execScript(elem)
		}
		// no alerts if maxAlertCounter is reached
		if maxAlertCounter > elem.MaxAlerts {
			alert = false
		}
		// alert triggers
		if alert && elem.Email {
			sendEmail(elem, err)
		}
		if alert && elem.Log {
			sendLog(elem, err)
		}

		// reset the alert and sleep
		alert = false
		time.Sleep(time.Duration(elem.Interval) * time.Second)
	}
}


func main() {
	// TODO: json file name via argument
	// TODO: list of json files, via argument
	json_file := "./files/gocontrol.json"
	// get the requests to monitor

	res, err := GetRequestInfo(json_file)
	if  err != nil {
		log.Fatal("Fatal error loading the JSON data,", err)
	}
	for _, elem := range res {
		fmt.Printf("URL: %s Script:%s (Secure:%v) \n", elem.Url, elem.Script, elem.Secure)
	}

	// create go routines for each request
	var wg sync.WaitGroup
	for _, elem := range res {
		wg.Add(1)
		go startWorker(elem, &wg)
	}
	wg.Wait()

}

