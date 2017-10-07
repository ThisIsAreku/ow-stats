package main

import (
	"log"
	"encoding/json"
	"fmt"
	"os"
	ow "github.com/ThisIsAreku/ow-stats"
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
