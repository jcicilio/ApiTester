{
	"TestSuiteName": "Sample Tests",
	"Tests":[
		{
			"TestName":"01",
			"Uri":"",
			"Method":"POST",
			"Body":"",
			"Expects": {"ReturnCode":404, "MaxSeconds":0.300}
		},
		{
			"TestName":"02",
			"Uri":"",
			"Method":"GET",
			"Body":"",
			"Expects": {"ReturnCode":200, "MaxSeconds":0.1}
		},
		{
			"TestName":"03",
			"Uri":"api/log",
			"Method":"POST",
			"Body":"{\"a\":100, \"b\":200}",
			"Expects": {"ReturnCode":201, "MaxSeconds":0.1}
		}
	]
}