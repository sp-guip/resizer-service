package main

import (
	"fmt"
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
		fmt.Println(testCase)
	}
}
