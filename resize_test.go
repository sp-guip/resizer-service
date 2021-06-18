package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/disintegration/imaging"
)

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
		url:               "file:///OgSquare.png",
		verticalPadding:   0,
		horizontalPadding: 0.025,
		width:             200,
		height:            200,
		targetWidthRatio:  .65,
		targetHeightRatio: .75,
	},
	{
		url:               "file:///OgSquare.png",
		verticalPadding:   0.025,
		horizontalPadding: 0.05,
		width:             200,
		height:            200,
		targetWidthRatio:  1.2,
		targetHeightRatio: 1.1,
	},
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
	baseUrl = fmt.Sprintf("http://localhost:%s/thumbnail", port)
}

// Test the service with a valid set of data for the correct cropped image behavior
// Expect not to get an error
func TestResize(t *testing.T) {
	for i, testCase := range testCases {
		t.Logf("Iter#%d", i)
		var newWidth = int(float32(testCase.width) * testCase.targetWidthRatio)
		var newHeight = int(float32(testCase.height) * testCase.targetHeightRatio)
		var url = fmt.Sprintf("%s?url=%s&width=%d&height=%d", baseUrl, testCase.url, newWidth, newHeight)
		response, err := http.Get(url)
		if err != nil {
			t.Errorf("Error requesting URL: %s, error: %s", url, err.Error())
			continue
		}
		if response.StatusCode != http.StatusOK {
			var res, _ = ioutil.ReadAll(response.Body)
			t.Errorf("Unexpected error for URL: %s, errorStatus: %s, error: %s", url, response.Status, res)
			continue
		}
		img, err := imaging.Decode(response.Body)
		if err != nil {
			t.Errorf("Error reading image output of the resizer service, error: %s", err.Error())
			continue
		}
		if img.Bounds().Dx() != newWidth {
			t.Errorf("Wrong width returned from service, expected: %d, actual: %d", newWidth, img.Bounds().Dx())
		}
		if img.Bounds().Dy() != newHeight {
			t.Errorf("Wrong height returned from service, expected: %d, actual: %d", newHeight, img.Bounds().Dy())
		}
	}
}
