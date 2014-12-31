// ApiTester
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
)

// Base URL for tests, passed in for each program run
// Provides a way to point same tests to different server
// at runtime
var BaseUrl string

// TestSuite Configuration File Format
type TestSuite struct {
	TestSuiteName string
	Tests         []TestSuiteSetup
}

type TestSuiteSetup struct {
	TestName string
	Uri      string
	Method   string
	Body     string
	Expects  TestExpectation
}

type TestExpectation struct {
	ReturnCode int
}

func RunTestSuite(fileName string) error {
	//http://stackoverflow.com/questions/16681003/how-do-i-parse-a-json-file-into-a-struct-with-go
	var testSetup TestSuite

	// Wait group to finish when all tests are complete
	wg := new(sync.WaitGroup)

	// Open configuration file
	configFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println("RunTestSuite: error opening config file - ", fileName, " ", err.Error())
		return err
	}

	// Parse configuration
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&testSetup); err != nil {
		fmt.Println("RunTestSuite: error parsing config file - ", fileName, " ", err.Error())
		return err
	}

	fmt.Printf("%#v\n", testSetup)

	// How many tests to wait for
	wg.Add(len(testSetup.Tests))

	for _, v := range testSetup.Tests {
		go RunTest(v, wg)
	}

	// Wait for all tests to complete
	wg.Wait()
	return nil
}

func RunTest(test TestSuiteSetup, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("Running test:", test.TestName)
}

func main() {
	// Get command line args for Service base URL
	flag.StringVar(&BaseUrl, "url", "", "The base URL for services")
	flag.Parse()
	if BaseUrl == "" {
		log.Fatal("url parameter is required")
		return
	}

	fmt.Println("BaseUrl: ", BaseUrl)

	// Read the test definition file - TBD: support files from a directory
	// Format of filename: *.tapi.js
	// Currently supports: test.tapi.js
	_ = RunTestSuite("test.tapi.js")

	// Launch go routine for each test (concurrency)
}
