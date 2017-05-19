[![Build Status](https://travis-ci.org/pedromg/gocontrol.svg?branch=master)](https://travis-ci.org/pedromg/gocontrol)
gocontrol
=========

Very simple Go app to check if app(s) is(are) up, and relaunch it(them).
`statuscode` refers to the single status codes and not the family (500 does not serve 501,502,...)

### Configuration file format

Format file: JSON structured:
```json
[
	{"name": "Github index",
	"url": "github.com", 
	"secure": true,
	"interval": 60, 
	"statuscode": [401,500,501,502,503],
	"maxalerts": 2,
	"script": "./script_1.sh", 
	"delayedby": 2,
	"email": true, 
	"smtphost": "mailtrap.io",
	"smtpport":2525,
	"smtpemail": "sender@me.com",
	"smtpusername": "smtpusername123",
	"smtppassword": "123456",
	"emailaddress": "me@me.com", 
	"log": true, 
	"logfile": "gocontrol_1.log"},
	
	{"name": "Google index",
	"url": "google.com", 
	"secure": false,
	"interval": 60, 
	"statuscode": [404,501,502,503],
	"maxalerts": 2,
	"script": "./script_2.sh", 
	"delayedby": 2,
	"email": true,
	"smtphost": "mailtrap.io",
	"smtpport":2525,
	"smtpemail": "sender@me.com",
	"smtpusername": "smtpusername123",
	"smtppassword": "123456", 
	"emailaddress": "me@me.com", 
	"log": true, 
	"logfile": "gocontrol_2.log"}
]
```

### Use

When `maxalerts` (maximum number of alerts) is reached, the alerts will no longer be logged or sent. When the service is up again, returning no error, the maxalerts counter is reset to zero.

The `delayedby` int is the number of times the execution of the script is delayed by alert counts. After that it is executed. If 0, gets executed upon first alert.


### Configuration fields

-	__name__: (string) the name for this monitoring section.
-	__url__: (string) the address to monitor.
-	__secure__: (bool) http vs https.
-	__interval__: (int) the interval in seconds between monitor requests.
-	__statuscode__: ([]int) int array of the status codes to monitor and generate alert; name them all, not the family.
-	__maxalerts__: (int) the max number os alerts to be sent; after that it becomes silent until the service is up again.
-	__script__: (string) the script to run via `sh`
-	__delayedby__: (int) the number of detections before executing the script.
-	__email__: (bool) send an email ?
-	__smtphost__: (string) the hostname of the email provider
-	__smtpport__: (int) the port of the smtp host
-	__smtpemail__: (string) the email of the sender (from header)
-	__smtpusername__: (string) the username for the smtp auth
-	__smtppassword__: (string) the password for the smtp auth
-	__emailaddress__: (string) email to receive the alerts.
-	__log__: (bool) log ?
-	__logfile__: (string) # file to append the log.

### Cross Compile

If you are building on OSX for Linux usage, make sure your Go e prepared to generate binaries for other architectures. To enable it for Linux:

```
$ cd  $GOROOT/src
$ GOOS=linux GOARCH=386 ./make.bash
```
Then to generate a linux specific binary:
```
$ GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o gocontrol.linux gocontrol.go
```

