package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"gocv.io/x/gocv"
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
		var newWidth = int(float32(testCase.width) * testCase.targetWidthRatio)
		var newHeight = int(float32(testCase.height) * testCase.targetHeightRatio)
		var url = fmt.Sprintf("%s?url=%s&width=%d&height=%d", baseUrl, testCase.url, newWidth, newHeight)
		response, err := http.Get(url)
		if err != nil {
			t.Errorf("Error requesting URL: %s, error: %s", url, err.Error())
			continue
		}
		if response.StatusCode != http.StatusOK {
			t.Errorf("Unexpected error for URL: %s, errorStatus: %s", url, response.Status)
			continue
		}
		bytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Errorf("Error reading response from service error: %s", err.Error())
			continue
		}
		img, err := gocv.IMDecode(bytes, gocv.IMReadAnyColor)
		if err != nil {
			t.Errorf("Error reading image output of the resizer service, error: %s", err.Error())
			continue
		}
		if img.Cols() != newWidth {
			t.Errorf("Wrong width returned from service, expected: %d, actual: %d", newWidth, img.Cols())
		}
		if img.Rows() != newHeight {
			t.Errorf("Wrong height returned from service, expected: %d, actual: %d", newHeight, img.Rows())
		}
		if img.Channels() != 3 {
			t.Errorf("The service responded with the wrong number of channels: %d instead of 3", img.Channels())
		}
		if img.Type() != gocv.MatTypeCV8SC3 {
			img.ConvertTo(&img, gocv.MatTypeCV8SC3)
		}

		if byteScalarToInt(img.Sum()) < newWidth*newHeight*100 {
			t.Errorf("Expected a non-blank image")
		}
		fmt.Println(i, testCase)
	}
}

// Converts a scalar(B4) to an int
func byteScalarToInt(scalar gocv.Scalar) int {
	return int(scalar.Val1)<<24 | int(scalar.Val2)<<16 | int(scalar.Val3)<<8 | int(scalar.Val4)
}
