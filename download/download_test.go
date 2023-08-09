package download

import (
	"io"
	"os"
	"testing"
)

func MakeFakeRequest() io.Reader {

	f, _ := os.OpenFile("example.html", os.O_RDONLY, 0644)

	return f

}

/*func TestMakeFakeRequest(t *testing.T) {

	table, rowcount := ParseTable(MakeFakeRequest())

	fmt.Println(len(table), rowcount)

}*/

func TestExcel(t *testing.T) {

	table, _ := ParseTable(MakeFakeRequest())

	ProduceExcel("report_template.xlsx", table, "newfile.xlsx")

}
