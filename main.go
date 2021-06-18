package main

import (
	"fmt"
	"image"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"gocv.io/x/gocv"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", handleResizeImage)
	http.ListenAndServe(":"+port, nil)
}

func httpErrorF(res http.ResponseWriter, format string, args ...interface{}) {
	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte(fmt.Sprintf(format, args...)))
}

//Handle the images resize operations from
// /Get ?url=url&width=number&height=number
func handleResizeImage(res http.ResponseWriter, req *http.Request) {
	var err = req.ParseForm()
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("malformed request URL: " + req.URL.String()))
	}
	var widthParam = req.FormValue("width")
	var heightParam = req.FormValue("height")
	var urlParam = req.FormValue("url")
	var width, height int
	if widthParam == "" {
		httpErrorF(res, "malformed request URL, no width Parameter specified: %s", req.URL.String())

	} else if heightParam == "" {
		httpErrorF(res, "malformed request URL, no height Parameter specified: %s", req.URL.String())

	} else if urlParam == "" {
		httpErrorF(res, "malformed request URL, no url Parameter specified: %s", req.URL.String())

	} else if width, err = strconv.Atoi(widthParam); err != nil || width == 0 {
		httpErrorF(res, "malformed request URL, the width Parameter specified is illegal, URL: %s, Width: %s", req.URL.String(), widthParam)

	} else if height, err = strconv.Atoi(heightParam); err != nil || height == 0 {
		httpErrorF(res, "malformed request URL, the height Parameter specified is illegal, URL: %s, Height: %s", req.URL.String(), heightParam)

	} else {
		imgData, err := getImage(urlParam)
		if err != nil {
			httpErrorF(res, err.Error())
			return
		}
		var dstMat = gocv.NewMatWithSize(height, width, gocv.MatTypeCV8SC3)
		var widthResizeRatio = float64(width) / float64(imgData.Cols())
		var heightResizeRatio = float64(height) / float64(imgData.Rows())
		var minResizeRatio = widthResizeRatio
		if widthResizeRatio > heightResizeRatio {
			minResizeRatio = heightResizeRatio
		}
		if minResizeRatio > 1 {
			minResizeRatio = 1
		}
		var size = image.Point{X: int(minResizeRatio * float64(imgData.Cols())), Y: int(minResizeRatio * float64(imgData.Rows()))}
		var xOffset = (dstMat.Cols() - size.X) / 2
		var yOffset = (dstMat.Rows() - size.Y) / 2
		var cropOffset = image.Rect(xOffset, yOffset, dstMat.Cols()-xOffset, dstMat.Rows()-yOffset)
		gocv.Resize(*imgData, imgData, size, 0, 0, gocv.InterpolationNearestNeighbor)
		var dstROI = dstMat.Region(cropOffset)
		imgData.CopyTo(&dstROI)
		fmt.Println(imgData)
	}
	// os.Exit(0)
}

// Fetch and parse the image from the url
// Varifying it's validity in the end
func getImage(url string) (*gocv.Mat, error) {
	req, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("unable to get image from url: %s, err: %w", url, err)
	}
	encodedBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading image data from url: %s, err: %w", url, err)
	}
	img, err := gocv.IMDecode(encodedBytes, gocv.IMReadAnyColor)
	if err != nil || img.Empty() {
		return nil, fmt.Errorf("error decoding image data from url: %s, err: %w", url, err)
	}
	return &img, nil
}
