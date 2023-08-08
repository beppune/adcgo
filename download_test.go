package main

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func MakeFakeRequest() io.Reader {

	f, _ := os.OpenFile("example.html", os.O_RDONLY, 0644)

	return f

}

func TestMakeFakeRequest(t *testing.T) {

	table := ParseTable(MakeFakeRequest())

	fmt.Println(table)

}
