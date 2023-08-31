package kaba

import "github.com/xuri/excelize/v2"

func ParseKabaExcel(file *excelize.File, output chan<- *KabaEntry) {

	sheet := file.GetSheetList()[0]

	gcs := func(row int) *KabaEntry {
		cellname, _ := excelize.CoordinatesToCellName(1, row)
		tstamp, _ := file.GetCellValue(sheet, cellname)

		if tstamp == "" {
			return nil
		}

		cellname, _ = excelize.CoordinatesToCellName(2, row)
		sensor, _ := file.GetCellValue(sheet, cellname)

		cellname, _ = excelize.CoordinatesToCellName(3, row)
		text, _ := file.GetCellValue(sheet, cellname)

		return &KabaEntry{
			TimeStamp: tstamp,
			Support:   "",
			Name:      "",
			Sensor:    sensor,
			Text:      text,
		}
	}

	row := 1

	for v := gcs(row); v != nil; {

		output <- v

		row++

		v = gcs(row)

	}
	close(output)
}
