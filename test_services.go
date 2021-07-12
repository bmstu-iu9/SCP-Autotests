package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Test struct {
	TestCase []string `json:"TestCase"`
	Expected []string `json:"Expected"`
}

func getTestsFromFile(path string) ([]Test, error) {
	jsonStorage, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer jsonStorage.Close()

	data, err := ioutil.ReadAll(jsonStorage)
	if err != nil {
		return nil, err
	}

	var tests []Test

	if err = json.Unmarshal(data, &tests); err != nil {
		return nil, err
	}

	return tests, nil
}
