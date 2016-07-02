package httpget

import(
	"net/http"
	"net/url"
	"fmt"
	"io/ioutil"
	"strings"
)

func Httpreq(act string, uri string, post string, values url.Values) (body string){
	client := &http.Client{}
	
	req, err := http.NewRequest(act, uri, strings.NewReader(post))
	req.URL.RawQuery = values.Encode()
        //fmt.Printf("%s\n",req.URL.RawQuery)
        if err != nil {
                fmt.Println(err)
        }

        resp, _ := client.Do(req)
        defer resp.Body.Close()

        bodyb, _ := ioutil.ReadAll(resp.Body)
        //fmt.Println(string(bodyb))
	body = string(bodyb)
	
	return
}


func main() {
	client := &http.Client{}

	values := url.Values{}
	values.Add("title","ABC")
	values.Add("action","history")
	
	req, err := http.NewRequest("GET", "https://ja.wikipedia.org/w/index.php", nil)
	req.URL.RawQuery = values.Encode()
	//fmt.Printf("%s\n",string(values.Encode()))
	if err != nil {
		fmt.Println(err)
	}
	//defer req.Body.Close()
	
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(b))
}
