package xl

import (
	"archive/zip"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func printXml(st []byte, fileName string) {
	file, _ := os.Create(fileName)
	file.Write(st)
}

func printJson(st any, fileName string) {
	jn, _ := json.Marshal(st)
	file, _ := os.Create(fileName)
	file.Write(jn)
}

func LogJson(st any) {
	jn, _ := json.Marshal(st)
	println(string(jn))
}

func getLinksFromZipFile(f *zip.File) HLinksMap {
	fReader, err := f.Open()
	if err != nil {
		log.Fatalf("error opening file: %s", err.Error())
	}

	var sheet Sheet
	b, _ := ioutil.ReadAll(fReader)
	parseErr := xml.Unmarshal(b, &sheet)
	if parseErr != nil {
		log.Fatalf("error parsing bytes: %s", err.Error())
	}

	linksMap := HLinksMap{}
	for _, link := range sheet.Body.Links {
		linksMap[link.Ref] = link.Display
	}

	return linksMap
}

func findLinkSheet(zipFiles []*zip.File) *zip.File {
	var foundFile *zip.File
	for _, f := range zipFiles {
		if f.Name == LinkSheet {
			foundFile = f
			return foundFile
		}
	}

	return nil
}

func UnzipXL(path string, ch chan (HLinksMap)) {
	reader, err := zip.OpenReader(Path)
	if err != nil {
		fmt.Printf("error opening reader: %s", err.Error())
	}
	defer reader.Close()

	foundLinkSheet := findLinkSheet(reader.File)
	if foundLinkSheet == nil {
		log.Fatal("failed to find link sheet")
	}

	linksMap := getLinksFromZipFile(foundLinkSheet)

	ch <- linksMap
}

func makeCellKey(colNum int, rowNum int) string {
	return NumbersToColumnKey[colNum+1] + strconv.Itoa(rowNum+1)
}

func getRawRowsFromXL(path string) [][]string {
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Printf("error opening file: %s", err.Error())
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("error closing file: %s", err.Error())
		}
	}()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Printf("error getting rows: %s", err.Error())
	}

	return rows
}

func makeRow(row []string, rowNum int, headers []string) Row {
	retRow := Row{}
	for index, header := range headers {
		cellKey := makeCellKey(index, rowNum)
		retRow = append(retRow, Cell{Ref: cellKey, Value: row[index], Header: header})
	}

	return retRow
}

func ReadXL(path string, tableCh chan ([]Row), headersCh chan ([]string)) {
	rows := getRawRowsFromXL(path)
	var headers []string
	var table []Row

	for rowIndex, row := range rows {
		if rowIndex == 0 {
			headers = row
		} else {
			table = append(table, makeRow(row, rowIndex, headers))
		}
	}

	tableCh <- table
	headersCh <- headers
}

func combineResults(xlTable []Row, linksMap HLinksMap) []map[string]string {
	combined := []map[string]string{}
	for _, row := range xlTable {
		m := map[string]string{}
		for _, cell := range row {
			m[cell.Header] = cell.Value
			if value, ok := linksMap[cell.Ref]; ok {
				linkKey := "link." + cell.Header
				m[linkKey] = value
			}
		}
		combined = append(combined, m)
	}

	return combined
}

func ProcessXL(path string) []map[string]string {
	hCh := make(chan ([]string))
	xmlCh := make(chan (HLinksMap))
	xlCh := make(chan ([]Row))

	go func() {
		UnzipXL(path, xmlCh)
	}()
	go func() {
		ReadXL(path, xlCh, hCh)
	}()

	links := <-xmlCh
	xlData := <-xlCh

	return combineResults(xlData, links)
}
