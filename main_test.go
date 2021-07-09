package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestCheckEntry(t *testing.T) {
	jsonStorage, err := os.Open("tests/entry_point_tests.json")
	if err != nil {
		t.Error(err)
	}

	defer jsonStorage.Close()

	data, err := ioutil.ReadAll(jsonStorage)
	if err != nil {
		t.Error(err)
	}

	var tests []struct {
		TestPath string `json:"TestPath"`
		Expected  string `json:"Expected"`
	}

	if err = json.Unmarshal(data, &tests); err != nil {
		t.Error(err)
	}

	for _, test := range tests {
		t.Run(test.TestPath, func(t *testing.T){
			res, err := checkFile(test.TestPath)
			if err != nil {
				t.Error(err)
			} else {
				if test.Expected != strings.TrimSpace(string(res)) {
					t.Errorf("Incorrect result\nExpected: %s\nRecieved: %s", test.Expected, res)
				}
			}
		})
	}
}
