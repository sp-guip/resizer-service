package main

import (
	"fmt"
	"os"
	"testing"
)

//https://images.pexels.com/photos/658687/pexels-photo-658687.jpeg?cs=srgb&dl=pexels-cindy-gustafson-658687.jpg&fm=jpg
type testData struct {
	url               string
	verticalPadding   float32
	horizontalPadding float32
	width             int
	height            int
	targetWidthRatio  float32
	targetHeightRatio float32
}

var baseUrl string
var testCases = []testData{
	{
		url:               "https://images.pexels.com/photos/658687/pexels-photo-658687.jpeg?cs=srgb&dl=pexels-cindy-gustafson-658687.jpg&fm=jpg",
		verticalPadding:   0,
		horizontalPadding: 0,
		width:             2812,
		height:            2250,
		targetWidthRatio:  .8,
		targetHeightRatio: .8,
	},
	{
		url:               "https://images.pexels.com/photos/658687/pexels-photo-658687.jpeg?cs=srgb&dl=pexels-cindy-gustafson-658687.jpg&fm=jpg",
		verticalPadding:   0,
		horizontalPadding: 0.05,
		width:             2812,
		height:            2250,
		targetWidthRatio:  .7,
		targetHeightRatio: .8,
	},
	{
		url:               "https://images.pexels.com/photos/658687/pexels-photo-658687.jpeg?cs=srgb&dl=pexels-cindy-gustafson-658687.jpg&fm=jpg",
		verticalPadding:   0.1,
		horizontalPadding: 0,
		width:             2812,
		height:            2250,
		targetWidthRatio:  .7,
		targetHeightRatio: .5,
	},
	{
		url:               "https://images.pexels.com/photos/658687/pexels-photo-658687.jpeg?cs=srgb&dl=pexels-cindy-gustafson-658687.jpg&fm=jpg",
		verticalPadding:   0.05,
		horizontalPadding: 0.1,
		width:             2812,
		height:            2250,
		targetWidthRatio:  1.2,
		targetHeightRatio: 1.1,
	},
	{
		url:               "https://images.pexels.com/photos/658687/pexels-photo-658687.jpeg?cs=srgb&dl=pexels-cindy-gustafson-658687.jpg&fm=jpg",
		verticalPadding:   0.1,
		horizontalPadding: 0.1,
		width:             2812,
		height:            2250,
		targetWidthRatio:  1.2,
		targetHeightRatio: 1.2,
	},
}

func init() {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	baseUrl = fmt.Sprintf("http://localhost:%s/", port)
}

// Test the service with a valid set of data for the correct cropped image behavior
// Expect not to get an error
func TestResize(t *testing.T) {
	for i, testCase := range testCases {
		fmt.Println(i, testCase)
	}
}
