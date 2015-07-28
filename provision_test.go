package goonep

import ( 
	"testing"
	"github.com/stretchr/testify/assert"

	"strings"

    "encoding/json"
    "fmt"

)


func TestProvModel(t *testing.T) {

	rawData := "actived,abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde,{\u0026quot;A\u0026quot;:20,\u0026quot;B\u0026quot;:26,\u0026quot;C\u0026quot;:75,\u0026quot;PipeID\u0026quot;:\u0026quot;63mm\u0026quot;}"

	rawData = strings.Replace(rawData, "\u0026quot;", "\"",-1)

	provModel := ProvModel{}

	provModel.Parse( rawData )

	assert.Equal(t, true, provModel.Validate(), "Build ProvModel Failed: %v", provModel)
	assert.Equal(t, "manage/model", provModel.GetPath(), "ProvModel Path wrong: %v", provModel)

	dump(provModel)

}

func dump( o interface{}) {

    result, _ := json.Marshal(o)

    fmt.Printf("\n\n*****************************\nDump value: %s \n*****************************\n\n", string(result) )
}