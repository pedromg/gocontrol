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
	"emailaddress": "me@me.com", 
	"log": true, 
	"logfile": "gocontrol_2.log"}
]
```

### Use

When `maxalerts` (maximum number of alerts) is reached, the alerts will no longer be logged or sent. When the service is up again, returning no error, the maxalerts counter is reset to zero.

The `delayedby` int is the number of times the execution of the script is delayed by alert counts. After that it is executed. If 0, gets executed upon first alert.


### Configuration fields

```

	"name:" "github test" # the name for this monitoring section.
	"url": "domain.com/path/to/monitoring/controler" # the address to monitor.
	"secure": true # http vs https.
	"interval": 60 # the interval between monitor requests.
	"statuscode": [401,501,502,503] # the status codes to monitor and generate alert; name them all, not the family.
	"maxalerts": 2 # the max number os alerts to be sent; after that it becomes silent until the service is up again.
	"script": "./script_1.sh" # the script to run via `sh`
	"delayedby": 2 # the number of detections before executing the script.
	"email": true # send an email ?
	"emailaddress": "me@me.com" # email to receive the alerts.
	"log": true # log ?
	"logfile": "gocontrol_1.log" # file to append the log.
```

