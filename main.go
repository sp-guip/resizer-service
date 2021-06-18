package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
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
	}
	// os.Exit(0)
}
