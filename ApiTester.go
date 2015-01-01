// ApiTester
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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
	TestSuiteName         string
	TestSuiteDescription  string
	TestSuiteResultStatus bool
	TestSuiteTotalSeconds float64
	Tests                 []TestSuiteSetup
}

type TestSuiteSetup struct {
	TestName        string
	TestDescription string
	Uri             string
	Method          string
	Body            string
	Headers         []HeaderMap
	Expects         TestExpectation
	Result          TestResult
}

type TestExpectation struct {
	ReturnCode int
	MaxSeconds float64
	Headers    []HeaderMap
}

type TestResult struct {
	RunWhen              time.Time
	ReturnCode           int
	ReturnCodeStatusText string
	TestCompletionStatus bool
	Body                 string
	ElapsedTime          float64
}

type HeaderMap struct {
	Key   string
	Value string
}

func RunTestSuite(fileName string) (TestSuite, error) {
	//http://stackoverflow.com/questions/16681003/how-do-i-parse-a-json-file-into-a-struct-with-go
	var testSetup TestSuite

	// Open configuration file
	configFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println("RunTestSuite: error opening config file - ", fileName, " ", err.Error())
		return testSetup, err
	}

	// Parse configuration
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&testSetup); err != nil {
		fmt.Println("RunTestSuite: error parsing config file - ", fileName, " ", err.Error())
		return testSetup, err
	}

	//fmt.Println(testSetup.Tests[0].RequestHeaders)

	// Setup overall result status for future ANDing
	testSetup.TestSuiteResultStatus = true

	// Note that range returns copy of object, so use
	// the index if we want to modify the contents of base object
	for i, _ := range testSetup.Tests {
		RunTest(&testSetup.Tests[i])

		// Do some overall totals and overall result status
		testSetup.TestSuiteTotalSeconds += testSetup.Tests[i].Result.ElapsedTime
		testSetup.TestSuiteResultStatus = testSetup.TestSuiteResultStatus && testSetup.Tests[i].Result.TestCompletionStatus
	}

	// Output results to file
	if JsonOutputFile != "" {
		// Convert to nicely formatted json
		jsonByteArray, err := json.MarshalIndent(testSetup, "", "  ")
		if err != nil {
			fmt.Println("Error writing json formatted output file. ", err)
			return testSetup, err
		}

		f, err := os.Create(JsonOutputFile)
		if err != nil {
			fmt.Println("Error writing json formatted output file. ", err)
			return testSetup, err
		}

		defer f.Close()
		_, err = f.Write(jsonByteArray)
		if err != nil {
			fmt.Println("Error writing json formatted output file. ", err)
			return testSetup, err
		}
	}

	return testSetup, err
}

func RunTest(test *TestSuiteSetup) {
	fmt.Println("Running test:", test.TestName)
	var testResult TestResult

	// Call Uri
	// Setup request
	var jsonStr = []byte(test.Body)
	req, err := http.NewRequest(test.Method, BaseUrl+test.Uri, bytes.NewBuffer(jsonStr))
	if err != nil {
		testResult.TestCompletionStatus = false
		test.Result = testResult

		fmt.Println("Error setting up request for test: ", test.TestName, "  ", err)
		return
	}

	// Set the request headers
	for _, v := range test.Headers {
		req.Header.Set(v.Key, v.Value)
	}

	// Setup http client and start time
	start := time.Now()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		testResult.TestCompletionStatus = false
		test.Result = testResult

		fmt.Println("Error calling method ", err)
		return
	}
	defer resp.Body.Close()

	elapsed := time.Since(start)

	// Save the body of the returned API
	body, _ := ioutil.ReadAll(resp.Body)

	// Save test results and evaluate

	testResult.RunWhen = start
	testResult.ReturnCode = resp.StatusCode
	testResult.ReturnCodeStatusText = resp.Status
	testResult.ElapsedTime = elapsed.Seconds()
	testResult.Body = string(body)
	// Evaluate the headers: TBD

	testResult.TestCompletionStatus = (resp.StatusCode == test.Expects.ReturnCode) && (elapsed.Seconds() <= test.Expects.MaxSeconds)

	test.Result = testResult

	// Show consolidated run results
	fmt.Println("Response Headers:", resp.Header)
	fmt.Println("Response Body:", testResult.Body)
	fmt.Println("Response Status Text:", testResult.ReturnCodeStatusText)
	fmt.Println("Response Status Code:", test.Result.ReturnCode)
	fmt.Println("Expected Status Code:", test.Expects.ReturnCode)
	fmt.Println("Response Seconds: ", test.Result.ElapsedTime)
	fmt.Println("Expected Maximum Seconds: ", test.Expects.MaxSeconds)
	fmt.Println("Test Status:", test.Result.TestCompletionStatus, "\n")
}

func main() {
	// Get command line args for Service base URL
	flag.StringVar(&BaseUrl, "url", "", "The base URL for services")
	flag.StringVar(&JsonOutputFile, "json", "", "An optional filename, if supplied then test result and the test itself are output to json file.")
	flag.Parse()

	// Must have a base URL to run tests for
	if BaseUrl == "" {
		fmt.Println("url parameter is required")
		os.Exit(-1)
	}

	// Read the test definition file - TBD: support files from a directory
	// Format of filename: *.tapi.js
	// Currently supports: test.tapi.js
	result, err := RunTestSuite("test.tapi.js")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Println()
	fmt.Println("TestSuite Name: ", result.TestSuiteName)
	fmt.Println("TestSuite ElapsedTime (s): ", result.TestSuiteTotalSeconds)
	fmt.Println("TestSuite Overall Status: ", result.TestSuiteResultStatus)
	fmt.Println("TestSuite BaseUrl: ", BaseUrl)
	fmt.Println("TestSuite OutputFile: ", JsonOutputFile)
	fmt.Println()

	// If all tests pass then exit code zero
	if result.TestSuiteResultStatus {
		os.Exit(0)
	}

	os.Exit(-1)
}
