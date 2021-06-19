package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gocv.io/x/gocv"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server listening on port: " + port)
	http.HandleFunc("/thumbnail", handleResizeImage)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error initializing server, err: %s", err.Error())
	}
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
		logErrorf("Error formating a jsonError, error: %s", err.Error())
	}
	logErrorf(jsonErr.Error)
	res.Write(msg)
}

func logErrorf(err string, args ...interface{}) {
	log.SetOutput(os.Stderr)
	log.Printf(err, args...)
	log.SetOutput(os.Stdout)
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
		img, err := fetchImage(urlParam)
		if err != nil {
			httpErrorF(res, err.Error())
			return
		}
		output, err := resizeAndCrop(*img, width, height)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(fmt.Sprintf("error trying to convert the result image to bytes, err: %s", err.Error())))
		} else {
			res.WriteHeader(http.StatusOK)
			if _, err := res.Write(output.Bytes()); err != nil {
				log.Fatalf("error writing image response, err: %s", err.Error())
			}
		}
	}
}

// Transport independent process to resize and crop the image
// Returns the endcoded image bytes
func resizeAndCrop(img gocv.Mat, width, height int) (bytes.Buffer, error) {
	var dstMat = gocv.NewMatWithSize(height, width, gocv.MatTypeCV8SC3)
	var widthResizeRatio = float64(width) / float64(img.Cols())
	var heightResizeRatio = float64(height) / float64(img.Rows())
	var minResizeRatio = widthResizeRatio

	if widthResizeRatio > heightResizeRatio {
		minResizeRatio = heightResizeRatio
	}
	if minResizeRatio > 1 {
		minResizeRatio = 1
	}

	var size = image.Point{X: int(minResizeRatio * float64(img.Cols())), Y: int(minResizeRatio * float64(img.Rows()))}
	var xOffset = (dstMat.Cols() - size.X) / 2
	var yOffset = (dstMat.Rows() - size.Y) / 2
	var cropOffset = image.Rect(xOffset, yOffset, dstMat.Cols()-xOffset, dstMat.Rows()-yOffset)
	gocv.Resize(img, &img, size, 0, 0, gocv.InterpolationNearestNeighbor)
	// replacing, due to gocv bug
	// https://www.gitmemory.com/issue/hybridgroup/gocv/387/456821537
	// var dstROI = dstMat.Region(cropOffset)
	// imgData.CopyTo(&dstROI)
	if img.Type() != gocv.MatTypeCV8SC3 {
		img.ConvertTo(&img, gocv.MatTypeCV8SC3)
	}
	imgPixels, _ := img.DataPtrInt8()
	dstPixels, _ := dstMat.DataPtrInt8()
	for i := 0; i < size.Y; i++ {
		copy(dstPixels[((i+cropOffset.Min.Y)*width+cropOffset.Min.X)*3:((i+cropOffset.Min.Y)*width+cropOffset.Max.X)*3], imgPixels[i*size.X*3:(i+1)*size.X*3])
	}

	outputBytes, err := gocv.IMEncode(gocv.JPEGFileExt, dstMat)
	return *bytes.NewBuffer(outputBytes), err
}

// Fetch and parse the image from the url
// Verifying it's validity in the end
func fetchImage(url string) (*gocv.Mat, error) {
	var err error
	var input io.Reader
	if strings.Index(url, "file:///") == 0 {
		input, err = os.Open(url[8:])
	} else {
		var req *http.Response
		req, err = http.Get(url)
		input = req.Body
	}
	if err != nil {
		logErrorf("unable to get image from url: %s, err: %s", url, err.Error())
		return nil, fmt.Errorf("unable to get image from url: %s", url)
	}
	encodedBytes, err := ioutil.ReadAll(input)
	if err != nil {
		logErrorf("unable to get image from url: %s, err: %s", url, err.Error())
		return nil, fmt.Errorf("unable to get image from url: %s", url)
	}
	img, err := gocv.IMDecode(encodedBytes, gocv.IMReadAnyColor)
	if err != nil || img.Empty() {
		logErrorf("error decoding image data from url: %s, err: %s", url, err.Error())
		return nil, fmt.Errorf("error decoding image data from url: %s", url)
	}
	return &img, nil
}
