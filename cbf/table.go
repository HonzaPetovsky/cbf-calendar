package cbf

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

type Table struct {
	XMLName xml.Name `xml:"table"`
	Teams   []Team   `xml:"team"`
}

type Team struct {
	XMLName  xml.Name `xml:"team"`
	Id       string   `xml:"group"`
	Position string   `xml:"pos"`
}

func ImportTable(url string) (*Table, error) {
	xmlContent, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer xmlContent.Body.Close()

	byteValue, _ := ioutil.ReadAll(xmlContent.Body)

	var table Table
	xml.Unmarshal(byteValue, &table)

	return &table, nil
}

func FindPositionInTable(teamId string, table *Table) string {
	for i := 0; i < len(table.Teams); i++ {
		if table.Teams[i].Id == teamId {
			return table.Teams[i].Position
		}
	}
	return "0"
}
