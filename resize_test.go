package main

import (
	"fmt"
	"image"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/disintegration/imaging"
	// "gocv.io/x/gocv"
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

		if !checkScalarAvg(img, newWidth*newHeight*10) {
			t.Errorf("Expected a non-blank image")
		}
		imaging.Save(img, fmt.Sprintf("test#%d.out.jpg", i))

		// Check what area is supposed to be cropped
		var topCroppedPixels = int(float32(testCase.height) * testCase.horizontalPadding)
		var leftCroppedPixels = int(float32(testCase.width) * testCase.verticalPadding)
		var bottomCropOffset = newHeight - int(float32(testCase.height)*testCase.horizontalPadding)
		var rightCropOffset = newWidth - int(float32(testCase.width)*testCase.verticalPadding)
		// Calculate the regions
		var topCroppedRegion = imaging.Crop(img, image.Rect(0, 0, newWidth, topCroppedPixels))
		var leftCroppedRegion = imaging.Crop(img, image.Rect(0, topCroppedPixels, leftCroppedPixels, bottomCropOffset))
		var rightCroppedRegion = imaging.Crop(img, image.Rect(rightCropOffset, topCroppedPixels, newWidth, bottomCropOffset))
		var bottomCroppedRegion = imaging.Crop(img, image.Rect(0, bottomCropOffset, newWidth, newHeight))
		// Check the cropped area is black accommodating the JPG lose of color accuracy
		if checkScalarAvg(topCroppedRegion, -1) {
			t.Errorf("Bad cropping, the top crop area has color")
		}
		if checkScalarAvg(leftCroppedRegion, -1) {
			t.Errorf("Bad cropping, the left crop area has color")
		}
		if checkScalarAvg(rightCroppedRegion, -1) {
			t.Errorf("Bad cropping, the right crop area has color")
		}
		if checkScalarAvg(bottomCroppedRegion, -1) {
			t.Errorf("Bad cropping, the bottom crop area has color")
		}
	}
}

// Calculates per channel sum of pixels
func sumOfImage(region image.Image) [3]uint32 {
	var sum [3]uint32
	for y := 0; y < region.Bounds().Dx(); y++ {
		for x := 0; x < region.Bounds().Dx(); x++ {
			r, g, b, a := region.At(x, y).RGBA()
			r = uint32((float64(r) / float64(a)) * 255)
			g = uint32((float64(g) / float64(a)) * 255)
			b = uint32((float64(r) / float64(b)) * 255)
			sum[0] += r
			sum[1] += g
			sum[2] += b
		}
	}
	return sum
}

// Check the scalar sum of all pixel data doesn't pass more than 1 per pixel in average
func checkScalarAvg(region image.Image, max int) bool {
	if max == -1 {
		max = region.Bounds().Dx() * region.Bounds().Dy() * 3
	}
	var umax = uint32(max)
	var sums = sumOfImage(region)
	return sums[0] > umax || sums[1] > umax || sums[2] > umax
}
