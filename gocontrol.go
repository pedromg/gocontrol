package main

import (
	"os"
	"io/ioutil"
	"log"
	"fmt"
	"sync"
	"time"
	"bytes"
	"os/exec"
	"strconv"
	"math/rand"
	"encoding/json"
	"net/http"
	"net/mail"
	"net/smtp"
	"github.com/pedromg/goEncoderBase64"
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
	SmtpHost string
	SmtpPort int
	SmtpEmail string
	SmtpUsername string
	SmtpPassword string
	EmailAddress string
	Log bool
	LogFile string
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

func messageLine(elem RequestInfo, r *http.Response, err error) (m string) {
	var theErr string
	var theStatus string
	if err != nil {
		theErr = err.Error()
	}
	if r != nil {
		theStatus = r.Status
	}
	m = time.Now().String() + " " + elem.Name + " - " + theStatus + " - " + theErr + "\n"
	return m
}

func send_email(elem RequestInfo, from, to, msg string) {
	f := mail.Address{from, from}
	t := mail.Address{to, to}
	// auth
	auth := smtp.PlainAuth("", elem.SmtpUsername, elem.SmtpPassword, elem.SmtpHost)
	err := smtp.SendMail(elem.SmtpHost+":"+strconv.Itoa(elem.SmtpPort), auth, f.Address, []string{t.Address}, []byte(msg))
	if err != nil {
		log.Print("SEND MAIL Error, ", err)
	}
}

func sendEmail(elem RequestInfo, r *http.Response, err error) (e error) {
	startTime := time.Now()
        header := make(map[string]string)
	header["From"] = elem.SmtpEmail
	header["To"] = elem.EmailAddress
	the_mesg_id := "<" + strconv.Itoa(rand.Intn(999999999)) + "__" +
		startTime.Format("2006-01-02T15:04:05.999999999Z07:00") +
		"==@"+elem.SmtpHost+">" 
	header["Message-id"] = the_mesg_id
	header["Date"] = startTime.Format("Mon, 02 Jan 2006 15:04:05 +0000")
	header["Subject"] = "goControl Alert for " + elem.Name
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"
	body := "++ goControl ALERT ++ \n\n "
	body += messageLine(elem, r, err)
	msg := ""
	for k, v := range header {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msg += "\r\n"
	msg += goEncoderBase64.Base64MimeEncoder(body)
	go send_email(elem, elem.SmtpEmail, elem.EmailAddress, msg)

	return e
}

func sendLog(elem RequestInfo, r *http.Response, err error) (e error) {
	f, e := os.OpenFile(elem.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if e == nil {
		defer f.Close()
		_, e = f.WriteString(messageLine(elem, r, err))
	}
	return e
}

func execScript(elem RequestInfo) (err error) {
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
			sendEmail(elem, r, err)
		}
		if alert && elem.Log {
			e := sendLog(elem, r, err)
			if e != nil {
				log.Print("ERROR Creating LogFile - ", e)
			}
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

