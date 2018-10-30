package cbf

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

type Games struct {
	XMLName xml.Name `xml:"games"`
	Games   []Game   `xml:"game"`
}

type Game struct {
	XMLName xml.Name `xml:"game"`
	Id      string   `xml:"id"`
	Date    string   `xml:"gdate"`
	Time    string   `xml:"gtime"`
	Place   string   `xml:"place"`
	Teams   []Teams  `xml:"team"`
	Result  Result   `xml:"result"`
}

type Teams struct {
	XMLName xml.Name `xml:"team"`
	Id      string   `xml:"id"`
	Name    string   `xml:"name"`
}

type Result struct {
	XMLName xml.Name `xml:"result"`
	Score   Score    `xml:"score"`
}

type Score struct {
	XMLName xml.Name `xml:"score"`
	A       string   `xml:"a"`
	B       string   `xml:"b"`
}

func ImportGames(url string) (*Games, error) {
	xmlContent, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer xmlContent.Body.Close()

	byteValue, _ := ioutil.ReadAll(xmlContent.Body)

	var games Games
	xml.Unmarshal(byteValue, &games)

	return &games, nil
}
