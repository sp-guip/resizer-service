package main

import (
	"fmt"
	"net/http"
	"testing"
)

type partialTestData struct {
	url    string
	width  int
	height int
}

// ?url=http&width=0&height=432
// ?url=&width=132&height=432
// ?url=http&width=123&height=0
// ?width=123&height=432
// ?url=http&height=432
// ?url=http&width=123
// ?(empty)
var partialTestCases = []partialTestData{
	{url: "https://images.pexels.com/photos/658687/pexels-photo-658687.jpeg?cs=srgb&dl=pexels-cindy-gustafson-658687.jpg&fm=jpg", height: 432},
	{width: 123, height: 432},
	{url: "https://images.pexels.com/photos/658687/pexels-photo-658687.jpeg?cs=srgb&dl=pexels-cindy-gustafson-658687.jpg&fm=jpg", width: 123},
	{url: "none", width: 123, height: 432},
	{url: "https://images.pexels.com/photos/658687/pexels-photo-658687.jpeg?cs=srgb&dl=pexels-cindy-gustafson-658687.jpg&fm=jpg", width: -1, height: 432},
	{url: "https://images.pexels.com/photos/658687/pexels-photo-658687.jpeg?cs=srgb&dl=pexels-cindy-gustafson-658687.jpg&fm=jpg", width: 123, height: -1},
	{url: "none", width: -1, height: -1},
}

// Test a url that is not malformed but doesn't contain all required data
// Expects an error always
func TestPartial(t *testing.T) {
	for _, testCase := range partialTestCases {
		var requestUrl = baseUrl + "?"
		var paramSuffix = ""
		if testCase.url != "none" {
			requestUrl += fmt.Sprintf("%s%s=%s", paramSuffix, "url", testCase.url)
			paramSuffix = "&"
		}
		if testCase.width != -1 {
			requestUrl += fmt.Sprintf("%s%s=%d", paramSuffix, "width", testCase.width)
			paramSuffix = "&"
		}
		if testCase.height != -1 {
			requestUrl += fmt.Sprintf("%s%s=%d", paramSuffix, "height", testCase.height)
		}
		if res, err := http.Get(requestUrl); err != nil {
			t.Errorf("Error requesting URL: %s, error: %s", requestUrl, err.Error())
		} else if res.StatusCode == http.StatusOK {
			t.Errorf("Expected an errorStatus from a request with url: %s", requestUrl)
		}
	}
}
