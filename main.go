package main

import (
	"flag"
	"log"
)

var (
	SCPVersion   *string
	refalVersion *string
	file         *string
	rsd          *bool
	saveTime     *bool
)

func parseCommandLineFlags() {
	SCPVersion = flag.String("scp", "1", "supercompiler version")
	refalVersion = flag.String("v", "default", "refal version")
	file = flag.String("path", "tests/main_tests.json", "path to tests")
	rsd = flag.Bool("rsd", false, "delete rsd files after")
	saveTime = flag.Bool("time", false, "save scp-runtime values")
	flag.Parse()
}

func main() {
	parseCommandLineFlags()

	tests, err := getTestsFromFile(*file, 1)
	if err != nil {
		log.Fatal(err)
	}

	if err = runTests(tests.autotests); err != nil {
		log.Fatal(err)
	}
}
