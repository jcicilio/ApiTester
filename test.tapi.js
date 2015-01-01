{
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