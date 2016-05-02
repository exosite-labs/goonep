// Go library for the OnePlatform Provisioning API
// http://docs.exosite.com/provision/
package goonep

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	// "net/http/httputil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var VendorToken = ""

var PROVISION_BASE = "/provision"
var PROVISION_ACTIVATE = PROVISION_BASE + "/activate"
var PROVISION_DOWNLOAD = PROVISION_BASE + "/download"
var PROVISION_MANAGE = PROVISION_BASE + "/manage"
var PROVISION_MANAGE_MODEL = PROVISION_MANAGE + "/model/"
var PROVISION_MANAGE_CONTENT = PROVISION_MANAGE + "/content/"
var PROVISION_REGISTER = PROVISION_BASE + "/register"

var Pool = struct {
	sync.RWMutex
	devices map[string]*ProvModel
}{devices: make(map[string]*ProvModel)}

type ProvContent struct{}

type ProvGroup struct{}

type ProvModel struct {
	RawData string

	ActiveStatus string
	Rid          string
	SN           string

	ExtraField string
	TimeStamp  int64

	managebycik       bool
	managebysharecode bool
	url               string
}

func (m *ProvModel) GetPath() string {
	return PROVISION_MANAGE_MODEL
}

// Find is a helper function for finding model with characteristics contained in string argument
func (m *ProvModel) Find(modelName, id string) (ProvModel, error) {
	// Create one to store if we can fetch it in
	fetchedModel := ProvModel{}

	// Check for a bad id.
	if len(id) == 0 {
		e := errors.New(fmt.Sprintf("Zero length ID"))
		return fetchedModel, e
	}

	// Create a key to use to index the map that consists of the id and the
	// vendor token.  We may have the same id in different domains (e.g.
	// production and developemnt) and we need unique pool item for each.
	key := id + VendorToken

	// Lock access to the map to see if we have the requested ID currently stored
	Pool.RLock()
	device, ok := Pool.devices[key]
	Pool.RUnlock()
	if ok {
		// Yes, return it
		return *device, nil
	}

	log.Printf("ID %s/%s not in cache, fetching", modelName, id)

	var headers = http.Header{}
	result, err := ProvCall(m.GetPath()+modelName+"/"+id, VendorToken, "", "GET", false, headers)

	// Check for a bad provision call
	if err != nil {
		e := errors.New(fmt.Sprintf("Error finding ID '%s': %v ", id, err))
		return fetchedModel, e
	}

	rawData := strings.Trim(string(result.([]uint8)), "\r\n")

	// Check if the call succeeded but didn't return data we can use
	switch rawData {
	case "HTTP/1.1 404 Not Found":
		fallthrough
	case "HTTP/1.1 412 Precondition Failed":
		e := errors.New(fmt.Sprintf("Unexpected result finding ID '%s': %v", id, rawData))
		return fetchedModel, e
	}

	// Everything looks good, create the device
	fetchedModel.Parse(rawData)
	fetchedModel.SN = id
	fetchedModel.TimeStamp = time.Now().Unix()

	// Write it to the map
	Pool.Lock()
	Pool.devices[key] = &fetchedModel
	Pool.Unlock()

	return fetchedModel, nil
}

func (m *ProvModel) Parse(RawData string) {

	if len(RawData) <= 0 {
		return
	}

	m.RawData = RawData

	extraFieldFetcher := regexp.MustCompile("([a-zA-Z0-9]+,){2}")
	m.ExtraField = strings.Trim(extraFieldFetcher.ReplaceAllString(RawData, ""), "\"")

	efSlices := strings.Split(RawData, ",")

	if len(efSlices) <= 2 {
		return
	}

	m.ActiveStatus = efSlices[0]
	m.Rid = efSlices[1]

}

func (m *ProvModel) Validate() bool {

	if len(m.Rid) != 40 {
		return false
	}

	return true

}

func (m *ProvModel) Bytes() []byte {
	return []byte(m.RawData)
}

type ProvShare struct{}

var Provision struct {
	Manage struct {
		Content ProvContent
		Group   ProvGroup
		Model   ProvModel
		Share   ProvShare
	}

	Admin struct {
		Auth ProvModel
	}

	Register ProvModel
}

type ProvRestModel interface {

	// GetPath retrive the URL path for each different models
	GetPath() string

	Create(attr *interface{}) (Response, error)

	Find(id string) (Response, error)
	All() (Response, error)

	Update(attr *interface{}) (Response, error)
	Delete(attr *interface{}) (Response, error)
}

// ProvCall is a helper function that carries out HTTP requests for Provisioning API calls
func ProvCall(path, key, data, method string, managebycik bool, extra_headers http.Header) (interface{}, error) {
	client := &http.Client{}

	var serverUrl = ""
	serverUrl = "https://m2.exosite.com"

	//fmt.Printf(serverUrl + path + "\n\n")
	req, _ := http.NewRequest(method, serverUrl+path, strings.NewReader(data))
	req.Header = extra_headers
	if managebycik {
		req.Header.Add("X-Exosite-CIK", key)
	} else {
		req.Header.Add("X-Exosite-Token", key)
	}
	if method == "POST" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	}
	req.Header.Add("Accept", "text/plain, text/csv, application/x-www-form-urlencoded")

	// uncomment to print request (for debugging)
	// reqdump, _ := httputil.DumpRequestOut(req, true)
	// fmt.Printf("\r\n\r\n" + string(reqdump) + "\r\n\r\n")

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

// content_create implements POST to /provision/manage/content/<MODEL>/
func Content_create(provModel ProvModel, key, model, contentid, meta string, protect bool) (interface{}, error) {
	var data = "id=" + contentid + "&meta=" + meta
	if protect != false {
		data = data + "&protected=true"
	}
	var path = PROVISION_MANAGE_CONTENT + model + "/"
	var headers = http.Header{}
	return ProvCall(path, key, data, "POST", provModel.managebycik, headers)
}

// content_download implements GET to /provision/download
func Content_download(provModel ProvModel, cik, vendor, model, contentid string) (interface{}, error) {
	var data = "vendor=" + vendor + "&model=" + model + "&id=" + contentid
	var headers = http.Header{}
	headers.Add("Accept", "*")
	return ProvCall(PROVISION_DOWNLOAD, cik, data, "GET", provModel.managebycik, headers)
}

// content_info implements GET to /provision/manage/content/<MODEL>/<CONTENT_ID>
// or GET to /provision/download
func Content_info(provModel ProvModel, key, model, contentid, vendor string) (interface{}, error) {
	var headers = http.Header{}
	if vendor == "" {
		var path = PROVISION_MANAGE_CONTENT + model + "/" + contentid
		return ProvCall(path, key, "", "GET", provModel.managebycik, headers)
	} else {
		var data = "vendor=" + vendor + "&model=" + model + "&info=true"
		return ProvCall(PROVISION_DOWNLOAD, key, data, "GET", provModel.managebycik, headers)
	}
}

// content_list implements GET to /provision/manage/content/<MODEL>/
func Content_list(provModel ProvModel, key, model string) (interface{}, error) {
	var path = PROVISION_MANAGE_CONTENT + model + "/"
	var headers = http.Header{}
	return ProvCall(path, key, "", "GET", provModel.managebycik, headers)
}

// content_remove implements DELETE to /provision/manage/content/<MODEL>/<CONTENT_ID>
func Content_remove(provModel ProvModel, key, model, contentid string) (interface{}, error) {
	var headers = http.Header{}
	var path = PROVISION_MANAGE_CONTENT + model + "/" + contentid
	return ProvCall(path, key, "", "DELETE", provModel.managebycik, headers)
}

// content_upload implements POST to /provision/manage/content/<MODEL>/<CONTENT_ID>
func Content_upload(provModel ProvModel, key, model, contentid, data, mimetype string) (interface{}, error) {
	var headers = http.Header{}
	headers.Add("Content-Type", mimetype)
	var path = PROVISION_MANAGE_CONTENT + model + "/" + contentid
	return ProvCall(path, key, data, "POST", provModel.managebycik, headers)
}

// model_create implements POST to /provision/manage/model/
func Model_create(provModel ProvModel, key, model, sharecode string, aliases, comments, historical bool) (interface{}, error) {
	var headers = http.Header{}
	var data = "model=" + model
	if provModel.managebysharecode {
		data = data + "&code=" + sharecode
	} else {
		data = data + "&rid=" + sharecode
	}
	if aliases == false {
		data = data + "&options[]=noaliases"
	}
	if comments == false {
		data = data + "&options[]=nocomments"
	}
	if historical == false {
		data = data + "&options[]=nohistorical"
	}
	return ProvCall(PROVISION_MANAGE_MODEL, key, data, "POST", provModel.managebycik, headers)
}

// model_info implements GET to provision/manage/model/<MODEL>
func Model_info(provModel ProvModel, key, model string) (interface{}, error) {
	var headers = http.Header{}
	return ProvCall(PROVISION_MANAGE_MODEL+model, key, "", "GET", provModel.managebycik, headers)
}

// model_list implements GET to /provision/manage/model/
func Model_list(provModel ProvModel, key string) (interface{}, error) {
	var headers = http.Header{}
	return ProvCall(PROVISION_MANAGE_MODEL, key, "", "GET", provModel.managebycik, headers)
}

// model_remove implements DELETE to /provision/manage/model/<MODEL>
func Model_remove(provModel ProvModel, key, model string) (interface{}, error) {
	var headers = http.Header{}
	var data = "delete=true&model=" + model + "&confirm=true"
	var path = PROVISION_MANAGE_MODEL + model
	return ProvCall(path, key, data, "DELETE", provModel.managebycik, headers)
}

// model_update implements PUT to /provision/manage/model/<MODEL>
func Model_update(provModel ProvModel, key, model, clonerid string, aliases, comments, historical bool) (interface{}, error) {
	var headers = http.Header{}
	var data = "rid=" + clonerid
	var path = PROVISION_MANAGE_MODEL + model
	return ProvCall(path, key, data, "PUT", provModel.managebycik, headers)
}

// serialnumber_activate implements POST to /provision/activate
func Serialnumber_activate(provModel ProvModel, model, serialnumber, vendor string) (interface{}, error) {
	var headers = http.Header{}
	var data = "vendor=" + vendor + "&model=" + model + "&sn=" + serialnumber
	return ProvCall(PROVISION_ACTIVATE, "", data, "POST", provModel.managebycik, headers)
}

// serialnumber_add implements POST to /provision/manage/model/<MODEL>/
func Serialnumber_add(provModel ProvModel, key, model, sn string) (interface{}, error) {
	var headers = http.Header{}
	var data = "add=true&sn=" + sn
	var path = PROVISION_MANAGE_MODEL + model + "/"
	return ProvCall(path, key, data, "POST", provModel.managebycik, headers)
}

// serialnumber_add_batch implements POST to /provision/manage/model/<MODEL>/
func Serialnumber_add_batch(provModel ProvModel, key, model string, sns []string) (interface{}, error) {
	var headers = http.Header{}
	var data = "add=true"
	for i := range sns {
		data = data + "&sn[]=" + sns[i]
	}
	var path = PROVISION_MANAGE_MODEL + model + "/"
	return ProvCall(path, key, data, "POST", provModel.managebycik, headers)
}

// serialnumber_disable implements POST to /provision/manage/model/<MODEL>/<SN>
func Serialnumber_disable(provModel ProvModel, key, model, serialnumber string) (interface{}, error) {
	var headers = http.Header{}
	var data = "disable=true"
	var path = PROVISION_MANAGE_MODEL + model + "/" + serialnumber
	return ProvCall(path, key, data, "POST", provModel.managebycik, headers)
}

// serialnumber_enable implements POST to /provision/manage/model/<MODEL>/<SN>
func Serialnumber_enable(provModel ProvModel, key, model, serialnumber, owner string) (interface{}, error) {
	var headers = http.Header{}
	var data = "enable=true&owner=" + owner
	var path = PROVISION_MANAGE_MODEL + model + "/" + serialnumber
	return ProvCall(path, key, data, "POST", provModel.managebycik, headers)
}

// serialnumber_info implements GET to /provision/manage/model/<MODEL>/<SN>
func Serialnumber_info(provModel ProvModel, key, model, serialnumber string) (interface{}, error) {
	var headers = http.Header{}
	var path = PROVISION_MANAGE_MODEL + model + "/" + serialnumber
	return ProvCall(path, key, "", "GET", provModel.managebycik, headers)
}

// serialnumber_list implements GET to /provision/manage/model/<MODEL>/
func Serialnumber_list(provModel ProvModel, key, model string, offset, limit int) (interface{}, error) {
	var headers = http.Header{}
	var data = "offset=" + strconv.Itoa(offset) + "&limit=" + strconv.Itoa(limit)
	var path = PROVISION_MANAGE_MODEL + model + "/"
	return ProvCall(path, key, data, "GET", provModel.managebycik, headers)
}

// serialnumber_reenable implements POST to /provision/manage/model/<MODEL>/<SN>
func Serialnumber_reenable(provModel ProvModel, key, model, serialnumber string) (interface{}, error) {
	var headers = http.Header{}
	var data = "enable=true"
	var path = PROVISION_MANAGE_MODEL + model + "/" + serialnumber
	return ProvCall(path, key, data, "POST", provModel.managebycik, headers)
}

// serialnumber_remap implements POST to /provision/manage/model/<MODEL>/<SN>
func Serialnumber_remap(provModel ProvModel, key, model, serialnumber, oldsn string) (interface{}, error) {
	var headers = http.Header{}
	var data = "enable=true&oldsn=" + oldsn
	var path = PROVISION_MANAGE_MODEL + model + "/" + serialnumber
	return ProvCall(path, key, data, "POST", provModel.managebycik, headers)
}

// serialnumber_remove implements DELETE to /provision/manage/model/<MODEL>/<SN>
func Serialnumber_remove(provModel ProvModel, key, model, serialnumber string) (interface{}, error) {
	var headers = http.Header{}
	var path = PROVISION_MANAGE_MODEL + model + "/" + serialnumber
	return ProvCall(path, key, "", "DELETE", provModel.managebycik, headers)
}

// serialnumber_remove_batch implements POST to /provision/manage/model/<MODEL>/
func Serialnumber_remove_batch(provModel ProvModel, key, model string, sns []string) (interface{}, error) {
	var headers = http.Header{}
	var data = "remove=true"
	for i := range sns {
		data = data + "&sn[]=" + sns[i]
	}
	var path = PROVISION_MANAGE_MODEL + model + "/"
	return ProvCall(path, key, data, "POST", provModel.managebycik, headers)
}

// vendor_register implements POST to /provision/register
func Vendor_register(provModel ProvModel, key, vendor string) (interface{}, error) {
	var headers = http.Header{}
	var data = "vendor=" + vendor
	return ProvCall(PROVISION_REGISTER, key, data, "POST", provModel.managebycik, headers)
}

// vendor_show implements GET to /provision/register
func Vendor_show(key string) (interface{}, error) {
	var headers = http.Header{}
	return ProvCall(PROVISION_REGISTER, key, "", "GET", false, headers)
}

// vendor_unregister implements POST to /provision/register
func Vendor_unregister(key, vendor string) (interface{}, error) {
	var headers = http.Header{}
	var data = "delete=true&vendor=" + vendor
	return ProvCall(PROVISION_REGISTER, key, data, "POST", false, headers)
}
