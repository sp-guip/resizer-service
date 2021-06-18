package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	// "gocv.io/x/gocv"
	"github.com/disintegration/imaging"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/thumbnail", handleResizeImage)
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
		var originalWidth = imgData.Bounds().Dx()
		var originalHeight = imgData.Bounds().Dy()

		var widthResizeRatio = float64(width) / float64(originalWidth)
		var heightResizeRatio = float64(height) / float64(originalHeight)
		var minResizeRatio = math.Min(math.Min(widthResizeRatio, heightResizeRatio), 1)

		var size = image.Point{X: int(minResizeRatio * float64(originalWidth)), Y: int(minResizeRatio * float64(originalHeight))}
		var resizedImage = imaging.Resize(imgData, size.X, size.Y, imaging.BSpline)

		var dstImage = imaging.New(width, height, color.Black)
		var xOffset = (width - size.X) / 2
		var yOffset = (height - size.Y) / 2
		dstImage = imaging.Paste(dstImage, resizedImage, image.Point{X: xOffset, Y: yOffset})

		var buf = bytes.Buffer{}
		err = imaging.Encode(&buf, dstImage, imaging.JPEG, imaging.JPEGQuality(100))
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(fmt.Sprintf("error trying to convert the result image to bytes, err: %s", err.Error())))
		} else {
			res.WriteHeader(http.StatusOK)
			if _, err := res.Write(buf.Bytes()); err != nil {
				log.Fatalf("error writing image response, err: %s", err.Error())
			}
		}
	}
}

// Fetch and parse the image from the url
// Varifying it's validity in the end
func getImage(url string) (image.Image, error) {
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
		return nil, fmt.Errorf("unable to get image from url: %s, err: %w", url, err)
	}
	img, err := imaging.Decode(input)
	if err != nil {
		return nil, fmt.Errorf("error decoding image data from url: %s, err: %w", url, err)
	}
	return img, nil
}
