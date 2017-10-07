package ow

import (
	"testing"
	"encoding/json"
	"github.com/go-test/deep"
	"github.com/parnurzeal/gorequest"
	"fmt"
	"os"
)

type RegionalProfile struct {
	Eu *Profile `json:"eu"`
}

func TestProfileParser_Parse(t *testing.T) {
	referenceData := RegionalProfile{}
	_, b, errs := gorequest.New().Get(fmt.Sprintf("https://owapi.net/api/v3/u/%s/blob", os.Getenv("BATTLETAG"))).End()
	if len(errs) != 0 {
		panic(errs[0].Error())
	}

	err := json.Unmarshal([]byte(b), &referenceData)
	if err != nil {
		panic(err.Error())
	}

	p := NewProfileParser("eu", os.Getenv("BATTLETAG"))
	testData, err := p.Parse()
	if err != nil {
		panic(err.Error())
	}

	if diff := deep.Equal(RegionalProfile{testData}, referenceData); diff != nil {
		for _, err := range diff {
			t.Error(err)
		}
	}
}
