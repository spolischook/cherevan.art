package googleDrive

import (
	"encoding/csv"
	"github.com/cherevan.art/src/artWork"
	"github.com/kpango/glg"
	"io"
	"os"
)

const inventoryFileID = "1MPkRsuhJhJdYlTCOxHFiJoFX3DZalZPbTOUObacJPg0"
const inventoryCsvFile = "inventory.csv"

var order = 0

func (s Service) MustGetArtWorks() artWork.ArtWorks {
	err := s.DownloadSheetAsCSV(inventoryFileID, inventoryCsvFile)
	if err != nil {
		glg.Fatalf("Error downloading sheet as csv: %v", err)
		return nil
	}

	// Open the file
	csvfile, err := os.Open(inventoryCsvFile)
	if err != nil {
		glg.Fatalf("Error opening file: %v", err)
		return nil
	}

	// Parse the file
	r := csv.NewReader(csvfile)
	headers, err := r.Read()
	csv := artWork.NewCsv(headers)

	var artWorks []artWork.ArtWork
	// Iterate through the records
	for {
		order++
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			glg.Fatalf("Error reading csv: %v", err)
		}

		artWork := csv.ParseRow(record)
		artWork.Order = order
		artWorks = append(artWorks, artWork)
	}

	return artWorks
}
