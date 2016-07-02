// Copyright (C) 2016 Makoto Imaizumi <Suisui@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// This program searchs wikidata items on list in textfile, 
// returns wikidata Q,Label,description
// which label prefix maches to search word 
// Usage) wikibasesearch filename

package main

import(
	"net/url"
	"fmt"
	"os"
	"bufio"
	"strings"
	"encoding/json"
	"./httpget"
	"strconv"
)

//wbsearchqueryの検索結果json
type SearchinfoT struct{
	Search string `json:"search"`
}

type MatchT struct{
	Type string     `json:"type"`
	Language string `json:"language"`
//	Text string     `json:"text"`
}

type SearchT struct{
//	Id string       `json:"id"`
//	Conurl string   `json:"concepturi"`
//	Url string      `json:"url"`
	Title string    `json:"title"`
//	Pageid json.Number      `json:"pageid,Number"`
	Label string    `json:"label"`
	Description string   `json:"description"`
	Match MatchT    `json:"match"`
	Aliases []string `json:"aliases"`
}

type SearchResultT struct{
	SearchInfo SearchinfoT `json:"searchinfo"`
	Search []SearchT `json:"search"`
	SearchContinue json.Number  `json:"search-continue,Number"`
	Success json.Number   `json:"success,Number"`
}

//出力用の構造体
type ResultT struct{
	title	string
	label	string
	desc	string
}



func WdSearchQuery(search string, format string, lang string, cont int)(resp string){
	values := url.Values{}
        values.Add("action","wbsearchentities")
        values.Add("search",search)
	values.Add("language", lang)
	values.Add("format",format)
	values.Add("limit", "50")
	values.Add("uselang",lang)
	if cont > 0 {
		values.Add("continue",strconv.Itoa(cont))
	}
//	values.Add("formatversion","2")


	resp = httpget.Httpreq("GET", "https://wikidata.org/w/api.php", "", values)
	
	//fmt.Printf("%s",resp)
	return
}

func WdSearchTransform(sr SearchResultT, result []ResultT) []ResultT {
	slen := len([]rune(sr.SearchInfo.Search))
	for _,v := range sr.Search {
		//fmt.Printf("%#+v\n", v)
		lbl := []rune(v.Label)
		// 前方一致のみを抽出
		if string(lbl[:slen]) == sr.SearchInfo.Search {
			var tmpR ResultT
			tmpR.title = v.Title
			tmpR.label = v.Label
			tmpR.desc = v.Description	
			result = append(result, tmpR)
		}
	}
	return result
}

func WdSearch(search string, lang string, limit int) string {
	if limit > 50 {
		limit = 50
	}
	cont := 0
	var result []ResultT
	for ;; {
		//fmt.Printf("req %d\n",cont)
		resp := WdSearchQuery(search,"json", "ja", cont)
		cont = 0	
		var sc SearchResultT
		if err := json.Unmarshal([]byte(resp), &sc); err != nil {
			panic(err)
		}

		result = WdSearchTransform(sc, result)

		contv, _ := json.Number.Int64(sc.SearchContinue)
		if contv != 0 {
			cont = int(contv)
		}else{
			break
		}
	}
	//1行分の文字列
	var linestr string
	linestr += search + "\t"
	for _, r := range result{
		//fmt.Printf("%#+v\n",r)
		linestr += fmt.Sprintf("%s\t%s\t%s\t",r.title, r.label, r.desc)
	}
	return linestr
}


func main(){
	var fp *os.File
	var err error
	if len(os.Args) < 2 {
		fp = os.Stdin
	} else {
		fp, err = os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		defer fp.Close()
	}

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		word := strings.Trim(scanner.Text(), "\t ")
		fmt.Printf("%s\n",WdSearch(word,"ja",153))
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

