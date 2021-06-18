package main

import (
	"fmt"
	"testing"
)

type malformedTestData struct {
	noQuerySign       bool
	partialParamNames bool
	url               string
	width             string
	height            string
	queryUrl          string
}

// URL=http&width=213&height=54
// ?UR=http&width=213&height=54
// ?URL=http:404&width=213&height=54
// ?URL=sad&width=213&height=54
// ?URL=http&width=a&height=54
// ?URL=http&width=14&height=b
// ?UR&L=http&width=14&height=b

var malformedTestCases = []malformedTestData{
	{
		noQuerySign: true,
		url:         "https://images.pexels.com/photos/658687/pexels-photo-658687.jpeg?cs=srgb&dl=pexels-cindy-gustafson-658687.jpg&fm=jpg",
		width:       "213",
		height:      "54",
	},
	{
		partialParamNames: true,
		url:               "https://images.pexels.com/photos/658687/pexels-photo-658687.jpeg?cs=srgb&dl=pexels-cindy-gustafson-658687.jpg&fm=jpg",
		width:             "213",
		height:            "54",
	},
	//404 page url
	{
		url:    "https://images.pexels.com/photos/658687/pexels-photo-658687.jpeg?cs=srgb&dl=pexels-cindy-gustafson-658687.jpg&fm=jpgdsa",
		width:  "213",
		height: "54",
	},
	//invalid page url
	{
		url:    "ht://www.pexels.com/photo/red-rose-flower-658687/",
		width:  "213",
		height: "54",
	},
	{
		url:    "https://www.pexels.com/photo/red-rose-flower-658687",
		width:  "a",
		height: "54",
	},
	{
		url:    "https://www.pexels.com/photo/red-rose-flower-658687",
		width:  "213",
		height: "b",
	},
	{
		url:    "https://www.pexels.com/photo/red-rose-flower-658687",
		width:  "a",
		height: "b",
	},
	{
		queryUrl: "?UR&L=https://images.pexels.com/photos/658687/pexels-photo-658687.jpeg?cs=srgb&dl=pexels-cindy-gustafson-658687.jpg&fm=jpg&width=54&height=214",
	},
}

// Tests badly formed urls that contain all the data required for the service but in an illegal format
// Expects an error, always
func TestMalformed(t *testing.T) {
	for _, testCase := range malformedTestCases {
		fmt.Println(testCase)
	}
}
