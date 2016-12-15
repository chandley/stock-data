package getFundamentals

import (
	//"errors"

	"net/http"
	"encoding/json"

	"fmt"
	"os"
)

type EdgarQuarterlyData struct {
	Result struct {
		       Totalrows int `json:"totalrows"`
		       Rows []struct {
			       Rownum int `json:"rownum"`
			       Values []struct {
				       Field string `json:"field"`
				       Value string `json:"value"`
			       } `json:"values"`
		       } `json:"rows"`
	       } `json:"result"`
}



func QuarterlyData() (err error, eqd *EdgarQuarterlyData) {
	url := generateUrl()
	if err != nil {
		return
	}
	fmt.Println(url)
	err = getJson(url, &eqd)
	fmt.Println(eqd)
	return
}

func generateUrl() (url string) {
	var apiKey string = os.Getenv("EDGAR_API_KEY")
	url = "http://edgaronline.api.mashery.com/v2/corefinancials/qtr.json?deleted=false&primarysymbols=msft&debug=false&sortby=primarysymbol+asc&appkey=" + apiKey
	return
}

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}



