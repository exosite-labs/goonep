package goonep

import (
	"math/rand"
	"runtime"
	"strconv"
	"testing"
	"time"
)

var cik = "<CIK>"
var alias = "X1"
var alias2 = "X2"

// errorCheckRPC checks for RPC API HTTP errors
func errorCheckRPC(t *testing.T, body Response, err interface{}, line int) {
	if err != nil {
		t.Errorf("Failed: %v", err)
	}
	if body.Results[0].Status == "invalid" {
		t.Errorf("Failed: %v", "RPC status response was invalid on line "+strconv.Itoa(line+1))
	}
	if body.Results[0].Status == "badarg" {
		t.Errorf("Failed: %v", "RPC status response was badarg on line "+strconv.Itoa(line+1))
	}
}

func TestMainRPC(t *testing.T) {
	var rid1 Response
	var rid2 Response
	var rid1Body interface{}
	var rid2Body interface{}
	var body Response
	var line int
	var desc = map[string]interface{}{
		"format":     "integer",
		"meta":       "",
		"name":       "who is me",
		"preprocess": []interface{}{},
		"public":     false,
		"retention": map[string]interface{}{
			"count":    "infinity",
			"duration": "infinity",
		},
		"subscribe": nil,
	}

	// lookup rid of alias X1, if doesn't exist then create + map
	rid1, err := Lookup(cik, "alias", alias)

	rid1Body = rid1.Results[0].Body
	if rid1.Results[0].Status == "invalid" {
		rid1, err = Create(cik, "dataport", desc)
		rid1Body = rid1.Results[0].Body
		_, _, line, _ = runtime.Caller(0)
		errorCheckRPC(t, rid1, err, line)
		body, err = OneMap(cik, rid1Body, alias)
		_, _, line, _ = runtime.Caller(0)
		errorCheckRPC(t, body, err, line)
	}

	// lookup rid of alias X2, if doesn't exist then create + map
	rid2, err = Lookup(cik, "alias", alias2)
	if err != nil {
		t.Errorf("Failed: %v", err)
	}
	rid2Body = rid2.Results[0].Body
	if rid2.Results[0].Status == "invalid" {
		rid2, err = Create(cik, "dataport", desc)
		rid2Body = rid2.Results[0].Body
		_, _, line, _ = runtime.Caller(0)
		errorCheckRPC(t, rid2, err, line)
		body, err = OneMap(cik, rid2Body, alias2)
		_, _, line, _ = runtime.Caller(0)
		errorCheckRPC(t, body, err, line)
	}

	// write data to dataport
	rand.Seed(time.Now().Unix())
	randomInt := rand.Intn(100-0) + 0
	body, err = Write(cik, rid1Body, randomInt)
	_, _, line, _ = runtime.Caller(0)
	errorCheckRPC(t, body, err, line)

	// read data from dataport
	body, err = Read(cik, rid1Body, map[string]interface{}{})
	_, _, line, _ = runtime.Caller(0)
	errorCheckRPC(t, body, err, line)

	// write group data
	time.Sleep(1 * time.Second)
	randomInt = rand.Intn(100-0) + 0
	array1 := []interface{}{rid1Body, randomInt}
	array2 := []interface{}{rid2Body, randomInt}
	entries := []interface{}{array1, array2}

	body, err = Writegroup(cik, entries)
	_, _, line, _ = runtime.Caller(0)
	errorCheckRPC(t, body, err, line)

	// read data from dataports
	body, err = Read(cik, rid1Body, map[string]interface{}{})
	_, _, line, _ = runtime.Caller(0)
	errorCheckRPC(t, body, err, line)

	body, err = Read(cik, rid2Body, map[string]interface{}{})
	_, _, line, _ = runtime.Caller(0)
	errorCheckRPC(t, body, err, line)

	// drop dataports
	body, err = Drop(cik, rid1Body)
	_, _, line, _ = runtime.Caller(0)
	errorCheckRPC(t, body, err, line)

	body, err = Drop(cik, rid2Body)
	_, _, line, _ = runtime.Caller(0)
	errorCheckRPC(t, body, err, line)

	// list client's dataports
	options := []interface{}{"dataport"}
	body, err = Listing(cik, options)
	_, _, line, _ = runtime.Caller(0)
	errorCheckRPC(t, body, err, line)

	// get all mapping aliases information of dataports
	// get resource id of device given key
	var device1 Response
	device1, err = Lookup(cik, "alias", "")
	_, _, line, _ = runtime.Caller(0)
	errorCheckRPC(t, device1, err, line)
	deviceRID := device1.Results[0].Body

	// get the alias information of given device
	option := map[string]interface{}{"aliases": true}
	body, err = Info(cik, deviceRID, option)
	_, _, line, _ = runtime.Caller(0)
	errorCheckRPC(t, body, err, line)
}
