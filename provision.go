// Go library for the OnePlatform Provisioning API
// https://github.com/exosite/docs/blob/master/provision/management.md
package goonep

import (
    "net/http"
    "log"
    "strings"
    "io/ioutil"
    "regexp"
    "time"
)

var VendorToken = ""

var Pool struct {
    Models  map[string]*ProvModel
}

type ProvContent struct {}

type ProvGroup struct {}

type ProvModel struct {
    RawData         string

    ActiveStatus    string
    Rid             string
    SN              string

    ExtraField      string
    TimeStamp       int64
}

func (m *ProvModel) GetPath() string {
    return "manage/model"
}

func (m *ProvModel) Find(modelName, id string) ProvModel {

    if Pool.Models[id] != nil {
        return *Pool.Models[id]
    }
    
    fetchedModel := ProvModel{}

    if len(id) <= 0 {
        log.Printf("Try find a non-sense ID: %d ", id)
        return ProvModel{}
    }

    result, err := ProvCall( VendorToken, "GET", m.GetPath() + "/" + modelName + "/" + id )

    if err != nil {
        log.Printf("Finding model(id: %s) met some error %v", id, err)
        return fetchedModel
    }

    rawData := strings.Trim( string( result.([]uint8) ), "\r\n")

    if rawData == "HTTP/1.1 404 Not Found" {
        return fetchedModel
    }

    fetchedModel.Parse(rawData)
    fetchedModel.SN = id
    fetchedModel.TimeStamp = time.Now().Unix()

    return fetchedModel
}

func (m *ProvModel) Parse( RawData string ) {

    if len( RawData ) <= 0 {
        return
    }

    m.RawData = RawData

    extraFieldFetcher := regexp.MustCompile( "([a-zA-Z0-9]+,){2}" )
    m.ExtraField = strings.Trim( extraFieldFetcher.ReplaceAllString( RawData, ""), "\"" )

    efSlices := strings.Split(RawData, "," )

    if len(efSlices) <= 2 {
        return
    }

    m.ActiveStatus = efSlices[0]
    m.Rid = efSlices[1]

}

func (m *ProvModel) Validate() bool {

    if len(m.Rid) != 40 { return false }

    return true

}

func (m *ProvModel) Bytes () []byte {
    return []byte( m.RawData )
}

type ProvShare struct {}

var Provision struct {

    Manage struct {
        Content ProvContent
        Group   ProvGroup
        Model   ProvModel
        Share   ProvShare
    }

    Admin struct {
        Auth    ProvModel
    }

    Register ProvModel

}

type ProvRestModel interface {

    // GetPath retrive the URL path for each different models
    GetPath () string

    Create  ( attr *interface{} ) Response

    Find    ( id string ) Response
    All     ( ) Response

    Update  ( attr *interface{} ) Response
    Delete  ( attr *interface{} ) Response

}

func ProvCall(vendorToken, method, path string) (interface{}, error) {
    client := &http.Client{}

    // https://m2.exosite.com/provision/manage/model/flow_sensor/POC_FLOW_01
    var serverUrl = ""
    serverUrl = "https://" + ONEPHost + "/provision/"
    if InDev {
        serverUrl = "https://m2-dev.exosite.com/provision/"
    }

    req, _ := http.NewRequest(method, serverUrl + path, nil)            
    req.Header.Add("X-Exosite-Token", vendorToken)

    resp, err := client.Do(req)

    if err != nil {
        return resp, err
    }

    defer resp.Body.Close()
    
    body, readErr := ioutil.ReadAll(resp.Body)
    if readErr != nil {
        return body, readErr
    }

    return body, nil
}
