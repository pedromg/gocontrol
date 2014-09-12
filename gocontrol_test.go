package main

import (
	"testing"
	"bytes"
	"os/exec"
	"github.com/pedromg/gocontrol"
)


const the_file = "./test/gocontrol.json"

func TestInexistentJSONFile( t *testing.T) {
	f := "./test/not_existing_agocontrol.json"
	_, e := main.GetRequestInfo(f)
	if e == nil {
		t.Error("JSON File Load Error:", e)
	}
}
func TestValidJSONFile( t *testing.T) {
	r, e := main.GetRequestInfo(the_file)
	if e != nil {
		t.Error("JSON File Load Error:", e)
	}
	if r[0].Url != "github.com" || r[1].Url != "google.com" {
		t.Error("Unmarshaled JSON error")
	}
	if r[0].Secure != true || r[1].Secure != false {
		t.Error("Unmarshaled JSON error")
	}
}

func TestScriptExists( t *testing.T) {
	r, _ := main.GetRequestInfo(the_file)
	for _, elem := range r {
		cmd := exec.Command("sh", elem.Script)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			t.Error("Script error,", err)
		}
	}
}

func TestValidRequest( t *testing.T) {
	t.Log("is this a valid check request ?")
}

func TestValidReply( t *testing.T) {
	t.Log("is the reply covered ?")
}

func TestSendEmail( t *testing.T) {
	t.Log("email sent ?")
}
