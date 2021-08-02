package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"
)

type Tests struct {
	unitTests []UnitTest
	autotests []Autotest
}

type Autotest struct {
	FilePath    string        `json:"FilePath"`
	TestData    string        `json:"TestData"`
	PerfectTime time.Duration `json:"PerfectTime"`
}

type UnitTest struct {
	TestCase []string `json:"TestCase"`
	Expected []string `json:"Expected"`
}

func getTestsFromFile(path string, num int) (*Tests, error) {
	jsonStorage, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer jsonStorage.Close()

	data, err := ioutil.ReadAll(jsonStorage)
	if err != nil {
		return nil, err
	}

	var tests Tests
	switch num {
	case 0:
		err = json.Unmarshal(data, &tests.unitTests)
	case 1:
		err = json.Unmarshal(data, &tests.autotests)
	default:
		err = errors.New("Only 2 types of tests supported\n")
	}
	if err != nil {
		return nil, err
	}

	return &tests, nil
}
