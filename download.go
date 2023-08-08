package main

import (
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/xuri/excelize/v2"
)

type AdcRequest http.Request

func panic_if_error(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func (req *AdcRequest) Prepare(rawurl string) *AdcRequest {
	u, err := url.Parse(rawurl)
	panic_if_error(err)

	r, err := http.NewRequest("POST", rawurl, nil)
	panic_if_error(err)

	r.Header.Add("Accept", `text/html`)
	r.Header.Add("Accept-Language", `it-IT,it;q=0.9`)
	r.Header.Add("Cache-Control", `max-age=0`)
	r.Header.Add("Connection", `keep-alive`)
	r.Header.Add("Content-Type", `application/x-www-form-urlencoded`)
	r.Header.Add("DNT", `1`)
	r.Header.Add("Origin", u.Hostname())
	r.Header.Add("Referer", rawurl)
	r.Header.Add("Upgrade-Insecure-Request", `1`)
	r.Header.Add("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36`)

	return req
}

func PrepareBody(s, bodyfile string) io.Reader {
	b, err := os.ReadFile(bodyfile)
	panic_if_error(err)

	b = append(b, s...)
	os.WriteFile("temp.txt", b, fs.ModeAppend)

	f, err := os.Open("temp.txt")
	panic_if_error(err)

	return f
}

func ParseTable(r io.Reader) ([][]string, int) {
	d, _ := goquery.NewDocumentFromReader(r)

	records := d.Find("#ADC_ContenutoSpecificoPagina_gvGiornaliero tr")

	numcols := records.Find("th").Length()
	cols := make([][]string, numcols)

	cells := records.Find("tr td")
	cellcount := cells.Length()

	for i := 0; i < cellcount; i++ {
		//fmt.Printf("Col num: %v, data: %v\n", i%numcols, strings.TrimSpace(cells.Slice(i, i+1).Text()))
		ar := &cols[i%numcols]
		*ar = append(*ar, cells.Slice(i, i+1).Text())
	}

	return cols, cellcount / numcols
}

func ProduceExcel(templatefile string, records [][]string, newfilename string) {
	file, err := excelize.OpenFile(templatefile, excelize.Options{})
	if err != nil {
		panic(err.Error())
	}

	sheet := file.GetSheetName(0)

	for col := 0; col < len(records); col++ {

		ar := records[col]

		for row := 0; row < len(ar); row++ {

			cell, _ := excelize.CoordinatesToCellName(col, row)
			file.SetCellValue(sheet, cell, strings.TrimSpace(ar[row]))

		}
	}

	file.SaveAs(newfilename)
}
