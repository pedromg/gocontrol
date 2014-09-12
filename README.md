gocontrol
=========

Very simple Go app to check if app(s) is(are) up, and relaunch it(them).

Format file: JSON structured:
```json
[
	{"url": "github.com", 
	"secure": true,
	"interval": 60, 
	"script": "./script_1.sh", 
	"email": true, 
	"email_address": "me@me.com", 
	"log": true, 
	"logfile": "gocontrol_1.log"},
	
	{"url": "google.com", 
	"secure": false,
	"interval": 60, 
	"script": "./script_2.sh", 
	"email": true, 
	"email_address": "me@me.com", 
	"log": true, 
	"logfile": "gocontrol_2.log"}
]
```

