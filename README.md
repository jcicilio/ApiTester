ApiTester
=========

ApiTesting Framework written in GO. 


# Current Expectations Supported

* Maximum Time To Complete
* Expected Headers
* Expected Return Code


# Usage

./ApiTester [Command Line Arguments]

Command Line Parameters

* --url = Service base URL to test, the URL of the API being tested
* --cfg = An optional test configuration file. Default to 'test.tapi.js'.
* --json = optional file name to write test result to
* --postpath = An optional url to the Api (see jcicilio/Api) to post test results to. Allows storing test results in MongoDb


# Test Configuration File

See --cfg Command Line Argument for name of file

Each test file defines a test suite.

* TestSuiteName - the name of the test suite
* TestSuiteDescription - a description of the test suite
* Tests - an array of test definitions
* Tests.TestName - the name of the test
* Tests.TestDescription - a description of what the test does
* Tests.Uri - the path from the root \ path of the route to test
* Tests.Method - the HTTP Method to use for the test
* Tests.Body - string encoded JSON
* Tests.Headers - an array of Key Value pairs defining the headers to use with the test
* Tests.Expects.ReturnCode - the expected return code for the test
* Tests.Expects.MaxSeconds - the maximum amount if time in seconds the test is expected to complete in
* Tests.Expects.Headers - an array of Key Value pairs defining the headers expected to be returned in the response


```
"TestSuiteName": "Sample Tests",
	"TestSuiteDescription":"A set of sample tests for evaluating the tester",
	"Tests":[
		{
			"TestName":"01",
			"TestDescription":"The test description 01",
			"Uri":"",
			"Method":"POST",
			"Body":"",
			"Headers":[{"Key":"Content-Type","Value":"application/json"},{"Key":"Authorization","Value":"SomeAuthKey"}],
			"Expects": {"ReturnCode":404, "MaxSeconds":0.300}
		},
		{
			"TestName":"02",
			"TestDescription":"The test description 02",
			"Uri":"",
			"Method":"GET",
			"Body":"",
			"Expects": {"ReturnCode":200, "MaxSeconds":0.10001,"Headers":[{"Key":"Content-type","Value":"application/json"}]}
		},
		{
			"TestName":"03",
			"TestDescription":"The test description 03",
			"Uri":"api/log",
			"Method":"POST",
			"Body":"{\"a\":100, \"b\":200}",
			"Expects": {"ReturnCode":201, "MaxSeconds":0.1}
		}
	]
}
```

## Features Considering

* Tests.Expects.BodyContains - expectation that the body of the response will contain specified contents