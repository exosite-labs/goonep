// Go library for the OnePlatform RPC
// http://docs.exosite.com/rpc/
package goonep

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	// "fmt"
)

var version = "0.2"
var DomainKey = ""

var InDev = false

// Set this to, e.g., "m2.exosite.com" or "localhost:18393"
var ONEPHost = "m2.exosite.com"

type Response struct {
	Results []Result
}

type Result struct {
	Id     int         `json:"id",int,omitempty`
	Body   interface{} `json:"result",string`
	Status string      `json:"status",string,omitempty`

	Error struct {
		Code    int
		Message string
	} `json:"error",omitempty`
}

// Call is a helper function that carries out HTTP requests for RPC API calls
func Call(auth string, procedure string, arguments []interface{}) (Response, error) {
	var calls = []interface{}{
		map[string]interface{}{
			"id":        1,
			"procedure": procedure,
			"arguments": arguments,
		},
	}
	return CallMulti(auth, calls)
}

func CallMulti(auth interface{}, calls []interface{}) (Response, error) {
	client := &http.Client{}

	f := Response{}

	var fullAuth = auth
	switch auth.(type) {
	case string:
		fullAuth = map[string]interface{}{"cik": auth}
	case interface{}:
		fullAuth = auth
	}

	var requestBody = map[string]interface{}{
		"auth":  fullAuth,
		"calls": calls,
	}

	var serverUrl = ""
	serverUrl = "https://" + ONEPHost + "/onep:v1/rpc/process"
	if InDev {
		serverUrl = "https://m2-dev.exosite.com/onep:v1/rpc/process"
	}

	buf, _ := json.Marshal(requestBody)
	requestBodyBuf := bytes.NewBuffer(buf)

	req, err := http.NewRequest("POST", serverUrl, requestBodyBuf)
	if err != nil {
		return f, err
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("User-Agent", "goonep "+version)

	resp, err := client.Do(req)

	if err != nil {
		return f, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return f, err
	}

	err2 := json.Unmarshal(body, &(f.Results))
	if err2 != nil {
		return f, err2
	}

	// uncomment to print response body (for debugging)
	// fmt.Println()
	// fmt.Println(string(body))
	// fmt.Println()
	// fmt.Println()
	// fmt.Println()

	// log.Printf("%v\n", f)
	// b, err := json.MarshalIndent(f, "", "  ")
	// log.Printf("%v\n", bytes.NewBuffer(b))

	// TODO: RPC error checking

	return f, nil
}

// the following functions implement the RPC APIs their names correspond to

func Activate(auth string, codetype string, code string) (Response, error) {
	var arguments = []interface{}{
		codetype,
		code,
	}
	return Call(auth, "activate", arguments)
}

func Create(auth string, ttype string, desc interface{}) (Response, error) {
	var arguments = []interface{}{
		ttype,
		desc,
	}
	return Call(auth, "create", arguments)
}

func Deactivate(auth string, codetype string, code string) (Response, error) {
	var arguments = []interface{}{
		codetype,
		code,
	}
	return Call(auth, "deactivate", arguments)
}

func Drop(auth string, rid interface{}) (Response, error) {
	var arguments = []interface{}{
		rid,
	}
	return Call(auth, "drop", arguments)
}

func Flush(auth string, rid interface{}) (Response, error) {
	var arguments = []interface{}{
		rid,
	}
	return Call(auth, "flush", arguments)
}

func Info(auth string, rid interface{}, options interface{}) (Response, error) {
	var arguments = []interface{}{
		rid,
		options,
	}
	return Call(auth, "info", arguments)
}

func Listing(auth string, types interface{}) (Response, error) {
	var arguments = []interface{}{
		types,
	}
	return Call(auth, "listing", arguments)
}

func Lookup(auth string, ttype string, alias string) (Response, error) {
	var arguments = []interface{}{
		ttype,
		alias,
	}
	return Call(auth, "lookup", arguments)
}

// oneMap implements the map RPC (name difference due to naming conflict)
func OneMap(auth string, rid interface{}, alias string) (Response, error) {
	var arguments = []interface{}{
		"alias",
		rid,
		alias,
	}
	return Call(auth, "map", arguments)
}

func Read(auth string, rid interface{}, options interface{}) (Response, error) {
	var arguments = []interface{}{
		rid,
		options,
	}
	return Call(auth, "read", arguments)
}

func Record(auth string, rid interface{}, entries interface{}, options interface{}) (Response, error) {
	var arguments = []interface{}{
		rid,
		entries,
		options,
	}
	return Call(auth, "record", arguments)
}

func Recordbatch(auth string, rid interface{}, entries interface{}) (Response, error) {
	var arguments = []interface{}{
		rid,
		entries,
	}
	return Call(auth, "recordbatch", arguments)
}

func Revoke(auth string, codetype string, code string) (Response, error) {
	var arguments = []interface{}{
		codetype,
		code,
	}
	return Call(auth, "revoke", arguments)
}

func Share(auth string, rid interface{}, options interface{}) (Response, error) {
	var arguments = []interface{}{
		rid,
		options,
	}
	return Call(auth, "share", arguments)
}

func Unmap(auth string, alias string) (Response, error) {
	var arguments = []interface{}{
		"alias",
		alias,
	}
	return Call(auth, "unmap", arguments)
}

func Update(auth string, rid interface{}, desc interface{}) (Response, error) {
	var arguments = []interface{}{
		rid,
		desc,
	}
	return Call(auth, "update", arguments)
}

func Usage(auth string, rid interface{}, metric string, starttime int, endtime string) (Response, error) {
	var arguments = []interface{}{
		rid,
		metric,
		starttime,
		endtime,
	}
	return Call(auth, "usage", arguments)
}

func Wait(auth string, rid interface{}) (Response, error) {
	var arguments = []interface{}{
		rid,
	}
	return Call(auth, "wait", arguments)
}

func Write(auth string, rid interface{}, value interface{}) (Response, error) {
	var arguments = []interface{}{
		rid,
		value,
	}
	return Call(auth, "write", arguments)
}

func Writegroup(auth string, entries interface{}) (Response, error) {
	var arguments = []interface{}{
		entries,
	}
	return Call(auth, "writegroup", arguments)
}
