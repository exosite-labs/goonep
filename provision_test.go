package goonep

import ( 
	"testing"
	// "github.com/stretchr/testify/assert"
	"strconv"
	"math/rand"
    "encoding/json"
    "fmt"
    "time"
    "runtime"
)

var vendorname = "VENDORNAMEHERE"
var vendortoken = "VENDORTOKENHERE"
var clonecik = "CLONECIKHERE"
var cloneportalcik = "CLONEPORTALCIKHERE" //use only if managing by sharecode
var portalcik = "PORTALCIKHERE"

/*func TestProvModel(t *testing.T) {

	rawData := "actived,abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde,{\u0026quot;A\u0026quot;:20,\u0026quot;B\u0026quot;:26,\u0026quot;C\u0026quot;:75,\u0026quot;PipeID\u0026quot;:\u0026quot;63mm\u0026quot;}"

	rawData = strings.Replace(rawData, "\u0026quot;", "\"",-1)

	provModel := ProvModel{}

	provModel.Parse( rawData )

	assert.Equal(t, true, provModel.Validate(), "Build ProvModel Failed: %v", provModel)
	assert.Equal(t, "manage/model", provModel.GetPath(), "ProvModel Path wrong: %v", provModel)

	dump(provModel)

}*/

func errorCheckProvision(t *testing.T, body string, err interface{}, line int) {
    if err != nil {
        t.Errorf("Failed: %v", err)
    }
    if body == "HTTP/1.1 409 Conflict\r\n" || body == "HTTP/1.1 404 Not Found\r\n" || body == "HTTP/1.1 412 Precondition Failed\r\n" {
        t.Errorf("Failed: %v", "HTTP status error on line " + strconv.Itoa(line+1))
    }
}

func TestMainProvision(t *testing.T) {
	rand.Seed(time.Now().Unix())
	randomInt := rand.Intn(10000000 - 0) + 0
	var model = "MyTestModel" + strconv.Itoa(randomInt)
	var sn1 = "001" + strconv.Itoa(randomInt)
	var sn2 = "002" + strconv.Itoa(randomInt)
	var sn3 = "003" + strconv.Itoa(randomInt)

	portalrid, err := lookup(portalcik, "alias", "")
	_, _, line, _ := runtime.Caller(0)
    errorCheckRPC(t, portalrid, err, line)
    portalridBody := portalrid.Results[0].Body
    fmt.Printf("portalrid: " + portalridBody.(string) + "\n\n")

    clonerid, err := lookup(clonecik, "alias", "")
    _, _, line, _ = runtime.Caller(0)
    errorCheckRPC(t, clonerid, err, line)
    cloneridBody := clonerid.Results[0].Body
    fmt.Printf("clonerid: " + cloneridBody.(string) + "\n\n")

    var meta = map[string]interface{}{
            "meta": "[\"" + vendorname + "\", \"" + model + "\"]",
    }
    sharecode, err := share(cloneportalcik, cloneridBody, meta)
    _, _, line, _ = runtime.Caller(0)
    errorCheckRPC(t, sharecode, err, line)
    sharecodeBody := sharecode.Results[0].Body

    provModel := ProvModel{
    	managebycik: false,
    	managebysharecode: true,
    	url: "http://m2.exosite.com",
    }
    body, err := model_create(provModel, vendortoken, model, sharecodeBody.(string), false, true, true)
    fmt.Printf(string(body.([]byte)) + "\r\n\r\n")
    _, _, line, _ = runtime.Caller(0)
    errorCheckProvision(t, string(body.([]byte)), err, line)

    body, err = model_list(provModel, vendortoken)
    fmt.Printf(string(body.([]byte)) + "\r\n\r\n")
    _, _, line, _ = runtime.Caller(0)
    errorCheckProvision(t, string(body.([]byte)), err, line)

    body, err = model_info(provModel, vendortoken, model)
    fmt.Printf(string(body.([]byte)) + "\r\n\r\n")
    _, _, line, _ = runtime.Caller(0)
    errorCheckProvision(t, string(body.([]byte)), err, line)

    body, err = serialnumber_add(provModel, vendortoken, model, sn1)
    fmt.Printf(string(body.([]byte)) + "\r\n\r\n")
    _, _, line, _ = runtime.Caller(0)
    errorCheckProvision(t, string(body.([]byte)), err, line)

   	var sn2andsn3 = []string{sn2, sn3}
   	body, err = serialnumber_add_batch(provModel, vendortoken, model, sn2andsn3)
   	fmt.Printf(string(body.([]byte)) + "\r\n\r\n")
    _, _, line, _ = runtime.Caller(0)
    errorCheckProvision(t, string(body.([]byte)), err, line)

   	body, err = serialnumber_list(provModel, vendortoken, model, 0, 10)
   	fmt.Printf(string(body.([]byte)) + "\r\n\r\n")
    _, _, line, _ = runtime.Caller(0)
    errorCheckProvision(t, string(body.([]byte)), err, line)

   	body, err = serialnumber_remove_batch(provModel, vendortoken, model, sn2andsn3)
   	fmt.Printf(string(body.([]byte)) + "\r\n\r\n")
    _, _, line, _ = runtime.Caller(0)
    errorCheckProvision(t, string(body.([]byte)), err, line)

   	body, err = serialnumber_list(provModel, vendortoken, model, 0, 1000)
   	fmt.Printf(string(body.([]byte)) + "\r\n\r\n")
    _, _, line, _ = runtime.Caller(0)
    errorCheckProvision(t, string(body.([]byte)), err, line)

   	fmt.Printf("serialnumber_enable() \r\n")
   	body, err = serialnumber_enable(provModel, vendortoken, model, sn1, portalridBody.(string))
   	fmt.Printf(string(body.([]byte)) + "\r\n\r\n")
    _, _, line, _ = runtime.Caller(0)
    errorCheckProvision(t, string(body.([]byte)), err, line)

   	body, err = serialnumber_info(provModel, vendortoken, model, sn1)
   	fmt.Printf(string(body.([]byte)) + "\r\n\r\n")
    _, _, line, _ = runtime.Caller(0)
    errorCheckProvision(t, string(body.([]byte)), err, line)

   	body, err = serialnumber_disable(provModel, vendortoken, model, sn1)
   	fmt.Printf(string(body.([]byte)) + "\r\n\r\n")
    _, _, line, _ = runtime.Caller(0)
    errorCheckProvision(t, string(body.([]byte)), err, line)

   	body, err = serialnumber_info(provModel, vendortoken, model, sn1)
   	fmt.Printf(string(body.([]byte)) + "\r\n\r\n")
    _, _, line, _ = runtime.Caller(0)
    errorCheckProvision(t, string(body.([]byte)), err, line)

   	body, err = serialnumber_reenable(provModel, vendortoken, model, sn1)
   	fmt.Printf(string(body.([]byte)) + "\r\n\r\n")
    _, _, line, _ = runtime.Caller(0)
    errorCheckProvision(t, string(body.([]byte)), err, line)

   	body, err = serialnumber_info(provModel, vendortoken, model, sn1)
   	fmt.Printf(string(body.([]byte)) + "\r\n\r\n")
    _, _, line, _ = runtime.Caller(0)
    errorCheckProvision(t, string(body.([]byte)), err, line)

   	body, err = serialnumber_activate(provModel, model, sn1, vendorname)
   	fmt.Printf(string(body.([]byte)) + "\r\n\r\n")
    _, _, line, _ = runtime.Caller(0)
    errorCheckProvision(t, string(body.([]byte)), err, line)

   	body, err = serialnumber_info(provModel, vendortoken, model, sn1)
   	fmt.Printf(string(body.([]byte)) + "\r\n\r\n")
    _, _, line, _ = runtime.Caller(0)
    errorCheckProvision(t, string(body.([]byte)), err, line)

   	body, err = content_create(provModel, vendortoken, model, "a.txt", "This is text", false)
   	fmt.Printf(string(body.([]byte)) + "\r\n\r\n")
    _, _, line, _ = runtime.Caller(0)
    errorCheckProvision(t, string(body.([]byte)), err, line)

   	body, err = content_upload(provModel, vendortoken, model, "a.txt", "This is content data", "text/plain")
   	fmt.Printf(string(body.([]byte)) + "\r\n\r\n")
    _, _, line, _ = runtime.Caller(0)
    errorCheckProvision(t, string(body.([]byte)), err, line)

   	body, err = content_list(provModel, vendortoken, model)
   	fmt.Printf(string(body.([]byte)) + "\r\n\r\n")
    _, _, line, _ = runtime.Caller(0)
    errorCheckProvision(t, string(body.([]byte)), err, line)

   	body, err = content_remove(provModel, vendortoken, model, "a.txt")
   	fmt.Printf(string(body.([]byte)) + "\r\n\r\n")
    _, _, line, _ = runtime.Caller(0)
    errorCheckProvision(t, string(body.([]byte)), err, line)

   	body, err = model_remove(provModel, vendortoken, model)
   	fmt.Printf(string(body.([]byte)) + "\r\n\r\n")
   	_, _, line, _ = runtime.Caller(0)
    errorCheckProvision(t, string(body.([]byte)), err, line)
}

func dump( o interface{}) {

    result, _ := json.Marshal(o)

    fmt.Printf("\n\n*****************************\nDump value: %s \n*****************************\n\n", string(result) )
}