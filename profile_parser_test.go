package ow_stats

import (
	"encoding/json"
	"fmt"
	"github.com/go-test/deep"
	"io/ioutil"
	"os"
	"testing"
)

type RegionalProfile struct {
	Eu *Profile `json:"eu"`
}

func TestProfileParser_Parse(t *testing.T) {
	referenceData := RegionalProfile{}
	r, err := Request.Get(fmt.Sprintf("https://owapi.net/api/v3/u/%s/blob", os.Getenv("BATTLETAG")))
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(b, &referenceData)
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
