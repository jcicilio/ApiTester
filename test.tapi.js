{
	"TestSuiteName": "Sample Tests",
	"TestSuiteDescription":"A set of sample tests for evaluating the tester",
	"Tests":[
		{
			"TestName":"01",
			"TestDescription":"Test of post to non-existant route",
			"Uri":"",
			"Method":"POST",
			"Body":"",
			"Headers":[{"Key":"Content-Type","Value":"application/json"},{"Key":"Authorization","Value":"SomeAuthKey"}],
			"Expects": {"ReturnCode":404, "MaxSeconds":0.300}
		},
		{
			"TestName":"02",
			"TestDescription":"Test of Get",
			"Uri":"",
			"Method":"GET",
			"Body":"",
			"IncludeBodyInResult": true,
			"Expects": {"ReturnCode":200, "MaxSeconds":0.10001,"Headers":[{"Key":"Content-type","Value":"application/json"}]}
		},
		{
			"TestName":"03",
			"TestDescription":"Test of Post",
			"Uri":"api/log",
			"Method":"POST",
			"Body":"{\"a\":100, \"b\":200}",
			"IncludeBodyInResult": false,
			"Headers":[{"Key":"Content-Type","Value":"application/json"},{"Key":"Authorization","Value":"SomeAuthKey"}],
			"Expects": {"ReturnCode":201, "MaxSeconds":0.1}
		},
		{
			"TestName":"04",
			"TestDescription":"Test of ApiSelector",
			"Uri":"apiSelector/circular/%7B%20%22storeId%22%3A%203403%7D",
			"Method":"GET",
			"Body":"{\"a\":100, \"b\":200}",
			"IncludeBodyInResult": false,
			"Headers":[{"Key":"Content-Type","Value":"application/json"}],
			"Expects": {"ReturnCode":200, "MaxSeconds":0.1}
		},
		{
			"TestName":"05",
			"TestDescription":"ApiSelector excluding body from result",
			"Uri":"apiSelector/circular/%7B%20%22storeId%22%3A%203403%7D",
			"Method":"GET",
			"Body":"{\"a\":100, \"b\":200}",
			"IncludeBodyInResult": true,
			"Headers":[{"Key":"Content-Type","Value":"application/json"}],
			"Expects": {"ReturnCode":200, "MaxSeconds":0.1}
		}
	]
}