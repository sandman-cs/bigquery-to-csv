package main

import (
	"encoding/json"
	"io/ioutil"
)

var (
	conf configuration
)

func init() {
	//Load Configuration Data
	dat, _ := ioutil.ReadFile("conf.json")
	err := json.Unmarshal(dat, &conf)
	checkFatalError("Could not load config file:", err)
}
