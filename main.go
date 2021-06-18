package main

import (
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", handleResizeImage)
	http.ListenAndServe(":"+port, nil)
}

//Handle the images resize operations from
// /Get ?url=url&width=number&height=number
func handleResizeImage(res http.ResponseWriter, req *http.Request) {
}
