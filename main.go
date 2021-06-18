package main

import (
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"log"
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

type jsonError struct {
	Error string
}

func httpErrorF(res http.ResponseWriter, format string, args ...interface{}) {
	res.WriteHeader(http.StatusBadRequest)
	var jsonErr = jsonError{
		Error: fmt.Sprintf(format, args...),
	}
	msg, err := json.Marshal(jsonErr)
	if err != nil {
		log.Fatalf("Error formating a jsonError, error: %s", err.Error())
	}
	res.Write(msg)
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
		// replacing, due to gocv bug
		// https://www.gitmemory.com/issue/hybridgroup/gocv/387/456821537
		// var dstROI = dstMat.Region(cropOffset)
		// imgData.CopyTo(&dstROI)
		if imgData.Type() != gocv.MatTypeCV8SC3 {
			imgData.ConvertTo(imgData, gocv.MatTypeCV8SC3)
		}
		imgPixels, _ := imgData.DataPtrInt8()
		dstPixels, _ := dstMat.DataPtrInt8()
		for i := 0; i < size.Y; i++ {
			copy(dstPixels[((i+cropOffset.Min.Y)*width+cropOffset.Min.X)*3:((i+cropOffset.Min.Y)*width+cropOffset.Max.X)*3], imgPixels[i*size.X*3:(i+1)*size.X*3])
		}

		outputBytes, err := gocv.IMEncode(gocv.JPEGFileExt, dstMat)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(fmt.Sprintf("error trying to convert the result image to bytes, err: %s", err.Error())))
		} else {
			res.WriteHeader(http.StatusOK)
			if _, err := res.Write(outputBytes); err != nil {
				log.Fatalf("error writing image response, err: %s", err.Error())
			}
		}
	}
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
