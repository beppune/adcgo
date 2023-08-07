package main

import (
	"io/fs"
	"io/ioutil"
	"os"
)

func main() {

	os.Create("temp.txt")

	//07%2F08%2F2023
	b, _ := ioutil.ReadFile("body.txt")

	b = append(b, `07%2F08%2F2023`...)
	ioutil.WriteFile("temp.txt", b, fs.ModeAppend)

}
