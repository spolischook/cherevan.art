package main

import (
	"encoding/csv"
	"github.com/cherevan.art/src/artWork"
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
	err = writer.Write(artWork.FbCsvHeaders())
	checkErr(err)

	// Write artworks data
	for _, aw := range existingArtWorks {
		awFb := aw.ToFb()
		err = writer.Write(awFb.ToCsv())
		checkErr(err)
	}
}
