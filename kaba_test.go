package main

import (
	"fmt"
	"testing"

	"github.com/beppune/adcgo/kaba"
	"github.com/xuri/excelize/v2"
)

func TestReportParser(t *testing.T) {

	in := make(chan *kaba.KabaEntry)

	file, err := excelize.OpenFile("2908_ALL.xlsx")
	if err != nil {
		panic(err)
	}

	go func() {
		kaba.ParseKabaExcel(file, in)
	}()

	for v := range in {
		fmt.Println(v)
	}
}
