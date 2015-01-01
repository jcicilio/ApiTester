// ApiTester
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// Base URL for tests, passed in for each program run
// Provides a way to point same tests to different server
// at runtime
var BaseUrl string
var JsonOutputFile string

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
	Result   TestResult
}

type TestExpectation struct {
	ReturnCode int
	MaxSeconds float64
}

type TestResult struct {
	RunWhen              time.Time
	ReturnCode           int
	TestCompletionStatus bool
	Body                 string
	ElapsedTime          float64
}

func RunTestSuite(fileName string) error {
	//http://stackoverflow.com/questions/16681003/how-do-i-parse-a-json-file-into-a-struct-with-go
	var testSetup TestSuite

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

	// Note that range returns copy of object, so use
	// the index if we want to modify the contents of base object
	for i, _ := range testSetup.Tests {
		RunTest(&testSetup.Tests[i])
	}

	// Output results to file
	if JsonOutputFile != "" {
		// Convert to nicely formatted json
		jsonByteArray, err := json.MarshalIndent(testSetup, "", "  ")
		if err != nil {
			fmt.Println("Error writing json formatted output file. ", err)
			return err
		}

		f, err := os.Create(JsonOutputFile)
		if err != nil {
			fmt.Println("Error writing json formatted output file. ", err)
			return err
		}

		defer f.Close()
		_, err = f.Write(jsonByteArray)
		if err != nil {
			fmt.Println("Error writing json formatted output file. ", err)
			return err
		}
	}

	return nil
}

func RunTest(test *TestSuiteSetup) {
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
		test.Result.TestCompletionStatus = false

		fmt.Println("Error calling method ", err)
		return
	}
	defer resp.Body.Close()

	elapsed := time.Since(start)

	// Save the body of the returned API
	body, _ := ioutil.ReadAll(resp.Body)

	var testResult TestResult
	// Save test results and evaluate
	testResult.RunWhen = start
	testResult.TestCompletionStatus = (resp.StatusCode == test.Expects.ReturnCode) && (elapsed.Seconds() < test.Expects.MaxSeconds)
	testResult.ReturnCode = resp.StatusCode
	testResult.ElapsedTime = elapsed.Seconds()
	testResult.Body = string(body)
	test.Result = testResult

	// Show consolidated run results
	fmt.Println("Response Headers:", resp.Header)
	//fmt.Println("Response Body:", string(body))
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Status Code:", test.Result.ReturnCode)
	fmt.Println("Expected Status Code:", test.Expects.ReturnCode)
	fmt.Println("Response Seconds: ", test.Result.ElapsedTime)
	fmt.Println("Expected Seconds: ", test.Expects.MaxSeconds)
	fmt.Println("Test Status:", test.Result.TestCompletionStatus, "\n")
}

func main() {
	// Get command line args for Service base URL
	flag.StringVar(&BaseUrl, "url", "", "The base URL for services")
	flag.StringVar(&JsonOutputFile, "json", "", "An optional filename, if supplied then test result and the test itself are output to json file.")
	flag.Parse()

	// Must have a base URL to run tests for
	if BaseUrl == "" {
		log.Fatal("url parameter is required")
		return
	}

	fmt.Println("BaseUrl       : ", BaseUrl)
	fmt.Println("JsonOutputFile: ", JsonOutputFile)
	fmt.Println()

	// Read the test definition file - TBD: support files from a directory
	// Format of filename: *.tapi.js
	// Currently supports: test.tapi.js
	_ = RunTestSuite("test.tapi.js")
}
