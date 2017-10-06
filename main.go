package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"encoding/json"
	"fmt"
	"github.com/ThisIsAreku/ow-stats/ow"
	"os"
)

const (
	URL = "https://playoverwatch.com/en-us/career/pc/eu/"
)

func fetchDocument(battleTag string) (doc *goquery.Document) {
	doc, err := goquery.NewDocument(URL + battleTag)
	if err != nil {
		panic(err.Error())
	}

	return doc
}

func main() {
	doc := fetchDocument(os.Getenv("BATTLETAG"))

	pp := ow.NewProfileParser(doc)

	profile, err := pp.Parse()
	if err != nil {
		log.Fatal(err.Error())
	}

	out, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(string(out))
}
