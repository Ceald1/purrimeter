package parsers
import (
	"strings"

	"github.com/sbabiv/xml2map"
)


func XML(xmlstr string) (out map[string]interface{}) {
	reader := strings.NewReader(xmlstr)
	decoder := xml2map.NewDecoder(reader)
	out, err := decoder.Decode()
	if err != nil {
		panic(err)
	}

	return
}