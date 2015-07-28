// Go library for the OnePlatform RPC
// https://github.com/exosite/api/tree/master/rpc
package goonep

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"
    //"log"
)

var version = "0.1"
var DomainKey = ""

var InDev = false
// Set this to, e.g., "m2.exosite.com" or "localhost:18393"
var ONEPHost = "m2.exosite.com"

type Response struct {
    Results []Result
}

type Result struct {
    Id      int         `json:"id",int,omitempty`
    Body    interface{} `json:"result",string`
    Status  string      `json:"status",string,omitempty`

    Error   struct {
        Code    int
        Message string
    } `json:"error",omitempty`
}


func Call(auth interface{}, procedure string, arguments []interface{}) (Response, error) {
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
    serverUrl = "https://" + ONEPHost + "/api:v1/rpc/process"
    if InDev {
        serverUrl = "https://m2-dev.exosite.com/api:v1/rpc/process"
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


    err2 := json.Unmarshal(body, &(f.Results) )
    if err2 != nil {
        return f, err2
    }

   // fmt.Printf("%v\n", f)
   // b, err := json.MarshalIndent(f, "", "  ")
   // fmt.Printf("%v\n", bytes.NewBuffer(b))

    // TODO: RPC error checking

    return f, nil
}
