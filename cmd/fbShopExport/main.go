package main

import (
	"encoding/csv"
	"github.com/cherevan.art/src/artWork"
	"github.com/cherevan.art/src/tool"
	"github.com/kpango/glg"
	"os"
)

func checkErr(err error) {
	if err != nil {
		glg.Fatal(err)
	}
}

func main() {
	existingArtWorks, err := artWork.GetExistingArtWorks()
	checkErr(err)

	file, err := os.Create("fbArtworks.csv")
	checkErr(err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	err = writer.Write(tool.CsvHeaders(artWork.Fb{}))
	checkErr(err)

	// Write artworks data
	for _, aw := range existingArtWorks {
		if aw.Validate() != nil {
			continue
		}
		if aw.Price < 1 {
			continue
		}
		awFb := aw.ToFb()
		err = writer.Write(tool.ToCsv(awFb))
		checkErr(err)
	}
}
