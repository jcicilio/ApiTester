/*

A testing framework for API testing.

Current expectations supported
	- Maximum Time To Complete

	- Expected Headers

	- Expected Return Code


Command Line Parameters

	--url = Service base URL to test, the URL of the API being tested

	--json = json file name to write result to, when provided a file to write test results to.

	--postpath = An optional url to the API to post test results to, allows for storing test results in MongoDb

*/
package main

import (
	"bytes"
	"encoding/json"
	"errors"
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
var PostPath string

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
	ErrorMessage         []string
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

	// Convert to nicely formatted json
	jsonByteArray, err := json.MarshalIndent(testSetup, "", "  ")
	if err != nil {
		fmt.Println("Error creating json for result output. ", err)
		return testSetup, err
	}

	// Output results to file if file path provided on command line
	if JsonOutputFile != "" {

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

	// Output results to API if api path provided on the command line
	if PostPath != "" {
		WriteTestSuiteResultToApi(jsonByteArray)
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

	// Evaluate the headers
	headersMatch, err := CheckHeaders(resp, test.Expects.Headers)
	if err != nil {
		testResult.ErrorMessage = append(testResult.ErrorMessage, fmt.Sprint(err))
	}

	// Evaluate the response code
	responseCodeMatch := (resp.StatusCode == test.Expects.ReturnCode)
	if !responseCodeMatch {
		testResult.ErrorMessage = append(testResult.ErrorMessage, fmt.Sprint("Response Code Mismatch"))
	}

	// Evaluate the response time
	expectedTimeMatch := (elapsed.Seconds() <= test.Expects.MaxSeconds)
	if !expectedTimeMatch {
		testResult.ErrorMessage = append(testResult.ErrorMessage, fmt.Sprint("Elapsed time greater than expected time"))
	}

	// Overall result of test
	testResult.TestCompletionStatus = responseCodeMatch && expectedTimeMatch && headersMatch

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

func CheckHeaders(resp *http.Response, expectedHeaders []HeaderMap) (bool, error) {
	// Check that each expected header matches the header response
	for _, h := range expectedHeaders {
		expectedHeader := resp.Header.Get(h.Key)
		if expectedHeader != h.Value {
			//fmt.Println("Header mismatch: ", h.Key, ":", h.Value, " != ", expectedHeader)
			return false, errors.New(fmt.Sprint("Header mismatch: ", h.Key, ":", h.Value, " != ", expectedHeader))
		}
	}

	return true, nil
}

func WriteTestSuiteResultToApi(testSuiteResults []byte) {
	// Setup for Posting Result
	req, err := http.NewRequest("POST", PostPath, bytes.NewBuffer(testSuiteResults))
	if err != nil {
		fmt.Println("Error setting up request for posting to API ", err)
		return
	}

	// Write the results to API
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		fmt.Println("Error writing test results to:  ", PostPath, " ", err)
		return
	}
}

func main() {
	// Get command line args for
	//
	// url = Service base URL to test
	// json = json file name to write result to
	// postpath = url of API to post test results to
	//
	flag.StringVar(&BaseUrl, "url", "", "The base URL for services")
	flag.StringVar(&JsonOutputFile, "json", "", "An optional filename, if supplied then test result and the test itself are output to json file.")
	flag.StringVar(&PostPath, "post", "", "An optional api route to post the test results to. Results are posted as json")
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
