package main

import (
	"net/http"
	"html/template"
	"encoding/xml"
	"io/ioutil"
)

type NewsMap struct {
	Keyword string
	Location string
}

type NewsAggPage struct {
    Title string
    News map[string]NewsMap
}

type Sitemapindex struct {
	Locations []string `xml:"sitemap>loc"`
}

type News struct {
	Titles []string `xml:"url>news>title"`
	Keywords []string `xml:"url>news>keywords"`
	Locations []string `xml:"url>loc"`
}



func newsAggHandler(w http.ResponseWriter, r *http.Request) {
	var s Sitemapindex
	var n News
	resp, _ := http.Get("https://www.washingtonpost.com/news-sitemap-index.xml")
	bytes, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes, &s)
	news_map := make(map[string]NewsMap)

	for _, Location := range s.Locations {
		resp, _ := http.Get(Location)
		bytes, _ := ioutil.ReadAll(resp.Body)
		xml.Unmarshal(bytes, &n)

		for idx, _ := range n.Keywords {
			news_map[n.Titles[idx]] = NewsMap{n.Keywords[idx], n.Locations[idx]}
		}
	}

	p := NewsAggPage{Title: "Go Aggregator", News: news_map}
    t, _ := template.ParseFiles("newsaggtemplate.html")
    t.Execute(w, p)
}

func main() {
	http.HandleFunc("/", newsAggHandler)
	http.ListenAndServe(":8000", nil) 
}