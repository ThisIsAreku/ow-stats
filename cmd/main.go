package main

import (
	"log"
	"os"
	ow "github.com/ThisIsAreku/ow-stats"
	"encoding/json"
	"fmt"
)

func main() {
	pp := ow.NewProfileParser("eu", os.Getenv("BATTLETAG"))

	profile, err := pp.Parse()
	if err != nil {
		log.Fatal(err.Error())
	}

	out, err := json.MarshalIndent(struct {
		Eu *ow.Profile `json:"eu"`
	}{
		profile,
	}, "", "  ")
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(string(out))
}
