package main

import (
	"bytes"
	"encoding/json"
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

	"github.com/disintegration/imaging"
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
		img, err := fetchImage(urlParam, res)
		if err != nil {
			httpErrorF(res, err.Error())
		}
		buf, err := resizeAndCrop(img, width, height)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			logErrorf(fmt.Sprintf("error trying to encode the image, err: %s", err.Error()))
			res.Write([]byte("error trying to encode the image"))
		} else {
			res.WriteHeader(http.StatusOK)
			if _, err := res.Write(buf.Bytes()); err != nil {
				logErrorf("error writing image response, err: %s", err.Error())
			}
		}
	}
}

func resizeAndCrop(img image.Image, width, height int) (bytes.Buffer, error) {
	var originalWidth = img.Bounds().Dx()
	var originalHeight = img.Bounds().Dy()

	var widthResizeRatio = float64(width) / float64(originalWidth)
	var heightResizeRatio = float64(height) / float64(originalHeight)
	var minResizeRatio = math.Min(math.Min(widthResizeRatio, heightResizeRatio), 1)

	var size = image.Point{X: int(minResizeRatio * float64(originalWidth)), Y: int(minResizeRatio * float64(originalHeight))}
	var resizedImage = imaging.Resize(img, size.X, size.Y, imaging.BSpline)

	var dstImage = imaging.New(width, height, color.Black)
	var xOffset = (width - size.X) / 2
	var yOffset = (height - size.Y) / 2
	dstImage = imaging.Paste(dstImage, resizedImage, image.Point{X: xOffset, Y: yOffset})

	var buf = bytes.Buffer{}
	return buf, imaging.Encode(&buf, dstImage, imaging.JPEG, imaging.JPEGQuality(100))
}

// Fetch and parse the image from the url
// Varifying it's validity in the end
func fetchImage(url string, res http.ResponseWriter) (image.Image, error) {
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
	img, err := imaging.Decode(input)
	if err != nil {
		logErrorf("error decoding image data from url: %s, err: %s", url, err.Error())
		return nil, fmt.Errorf("error decoding image data from url: %s", url)
	}
	return img, nil
}
