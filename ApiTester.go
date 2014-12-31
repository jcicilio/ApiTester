// ApiTester
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
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
	MaxSeconds float64
}

type TestResult struct {
	TestSetup            TestSuiteSetup
	ReturnCode           int
	TestCompletionStatus bool
}

func RunTestSuite(fileName string) error {
	//http://stackoverflow.com/questions/16681003/how-do-i-parse-a-json-file-into-a-struct-with-go
	var testSetup TestSuite

	// Wait group to finish when all tests are complete
	//wg := new(sync.WaitGroup)

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
	//wg.Add(len(testSetup.Tests))

	for _, v := range testSetup.Tests {
		//go RunTest(v)
		RunTest(v)
	}

	// Wait for all tests to complete
	//wg.Wait()
	return nil
}

func RunTest(test TestSuiteSetup) {
	//defer wg.Done()

	var testResult TestResult

	// Get test results
	testResult.TestSetup = test

	fmt.Println("Running test:", test.TestName)

	// Call Uri
	// Setup request
	var jsonStr = []byte(test.Body)
	req, err := http.NewRequest(test.Method, BaseUrl+test.Uri, bytes.NewBuffer(jsonStr))
	//req.Header.Set("X-Custom-Header", "myvalue")
	//req.Header.Set("Content-Type", "application/json")

	// Setup http client and start time
	start := time.Now()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		testResult.TestCompletionStatus = false

		fmt.Println("Error calling method ", err)
		return
	}
	defer resp.Body.Close()

	elapsed := time.Since(start)

	// Evaluate if test has passed
	testResult.TestCompletionStatus = (resp.StatusCode == test.Expects.ReturnCode) && (elapsed.Seconds() < test.Expects.MaxSeconds)

	//body, _ := ioutil.ReadAll(resp.Body)

	// Show consolidated run results
	fmt.Println("Response Headers:", resp.Header)
	//fmt.Println("Response Body:", string(body))
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Status Code:", resp.StatusCode)
	fmt.Println("Expected Status Code:", test.Expects.ReturnCode)
	fmt.Println("Response Seconds: ", elapsed.Seconds())
	fmt.Println("Expected Seconds: ", test.Expects.MaxSeconds)
	fmt.Println("Test Status:", testResult.TestCompletionStatus, "\n")
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
	// or maybe not - might be clearer to run sequentially
}
