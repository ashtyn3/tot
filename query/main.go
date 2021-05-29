package query

import (
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strings"

	"github.com/valyala/fastjson"
)

var Max *big.Int = big.NewInt(200000000)

func Req(path string, base bool) []byte {
	url := "https://en.wikipedia.org/api/rest_v1/page/random/summary"
	client := &http.Client{}

	req, _ := http.NewRequest("GET", url+path, nil)
	if base == false {
		req, _ = http.NewRequest("GET", path, nil)
	}
	resp, _ := client.Do(req)
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return body
}
func RandomRepo() string {
	var p fastjson.Parser
	treeBody := Req("", true)
	bTree, _ := p.ParseBytes(treeBody)
	contentRaw := bTree.Get("extract")
	if contentRaw == nil || len(contentRaw.String()) < 74 {
		return RandomRepo()
	}
	return contentRaw.String()
}

func formatUrl(str string, tempName string, new string) string {
	return strings.Replace(strings.Replace(str, "{/"+tempName+"}", new, 1), "\"", "", -1)
}
