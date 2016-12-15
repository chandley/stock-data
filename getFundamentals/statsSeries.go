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
				       Value forceToString `json:"value"`
			       } `json:"values"`
		       } `json:"rows"`
	       } `json:"result"`
}

type forceToString struct {
	content string
}

func (n *forceToString) UnmarshalJSON(buf []byte) error {
	n.content = string(buf)
	return nil
}



func QuarterlyData() (err error, eqd *EdgarQuarterlyData) {
	url := generateUrl("msft")
	if err != nil {
		return
	}
	fmt.Println(url)
	err = getJson(url, &eqd)
	fmt.Println(eqd)
	return
}

func LatestFundamentals(ticker string) (err error, data map[string]string) {
	url := generateUrl(ticker)
	if err != nil {
		return
	}
	eqd := new(EdgarQuarterlyData)
	err = getJson(url, &eqd)
	data = latestFundamentals(eqd)
	return
}

func generateUrl(ticker string) (url string) {
	var apiKey string = os.Getenv("EDGAR_API_KEY")
	url = "http://edgaronline.api.mashery.com/v2/corefinancials/qtr.json?deleted=false&primarysymbols=" + ticker + "&debug=false&sortby=primarysymbol+asc&appkey=" + apiKey
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

func latestFundamentals(eqd *EdgarQuarterlyData) map[string]string {
	fudamentals := make(map[string]string)

	latestDataSlice := eqd.Result.Rows[0].Values

	for _, pair := range latestDataSlice {
		fudamentals[pair.Field] = (pair.Value.content)
	}

	return fudamentals
}



