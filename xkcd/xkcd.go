package xkcd

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
)


// ----------------------------------------------------------------------------------
//  constants
// ----------------------------------------------------------------------------------

const (
	// URLs to interact with the xkcd api
	XKCD_URL = "http://xkcd.com/%d/info.0.json"
	XKCD_CURRENT_URL = "http://xkcd.com/info.0.json"

	// other constants
	CURRENT_COMIC = 0
)



// ----------------------------------------------------------------------------------
//  types
// ----------------------------------------------------------------------------------

type Comic struct {
	Alt        string `json:"alt"`
	Day        string `json:"day"`
	Img        string `json:"img"`
	Link       string `json:"link"`
	Month      string `json:"month"`
	News       string `json:"news"`
	Num        int    `json:"num"`
	SafeTitle  string `json:"safe_title"`
	Title      string `json:"title"`
	Transcript string `json:"transcript"`
	Year       string `json:"year"`
}


// ----------------------------------------------------------------------------------
//  functions
// ----------------------------------------------------------------------------------

func GetComic(num int) (*Comic, error) {
	// there is no commic with the number 0
	// so use this as the current comic
	url := fmt.Sprintf(XKCD_URL, num)
	if num == 0 {
		url = XKCD_CURRENT_URL
	}

	// do a http get request to fetch the 
	// json data for parsing
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// read the body of the request, containing the data
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// parse the json into an object
	var comic Comic
	err = json.Unmarshal(body, &comic)
	if err != nil {
		return nil, err
	}
	
	return &comic, nil
}
