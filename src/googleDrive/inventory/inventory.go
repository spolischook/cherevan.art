package inventory

import (
	"context"
	"encoding/csv"
	"github.com/cherevan.art/src/artWork"
	"github.com/cherevan.art/src/googleDrive/login"
	"github.com/kpango/glg"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"io"
	"os"
	"strconv"
)

const FileID = "1MPkRsuhJhJdYlTCOxHFiJoFX3DZalZPbTOUObacJPg0" // production sheet
//const FileID = "1gdLhmSfdJ3P2XddNjTUdjg7u7jQkErBrw19sTwYs_-I" // test sheet

type instance struct {
	sheets *sheets.Service
	drive  *drive.Service
}

func New() instance {
	tokenSource, err := login.TokenSource()

	// Use the TokenSource to authenticate the Sheets service
	s, err := sheets.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		glg.Fatal(err)
	}

	d, err := drive.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		glg.Fatal(err)
	}

	return instance{s, d}
}

func (i instance) SetId(rowNumber, id int) error {
	vr := &sheets.ValueRange{
		Values: [][]interface{}{{id}},
	}

	cellRange := "inventory!A" + strconv.Itoa(rowNumber+1)

	_, err := i.sheets.Spreadsheets.Values.
		Update(FileID, cellRange, vr).
		ValueInputOption("USER_ENTERED").
		Do()
	if err != nil {
		return err
	}

	return nil
}

func (i instance) GetArtWorks() artWork.ArtWorks {
	err := i.DownloadInventoryAsCSV(FileID, "inventory.csv")
	if err != nil {
		glg.Fatalf("Error downloading sheet as csv: %v", err)
		return nil
	}

	// Open the file
	csvFile, err := os.Open("inventory.csv")
	if err != nil {
		glg.Fatalf("Error opening file: %v", err)
		return nil
	}

	// Parse the file
	r := csv.NewReader(csvFile)
	headers, err := r.Read()
	artWorkCsv := artWork.NewCsv(headers)

	var order = 0
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

		w := artWorkCsv.ParseRow(record)
		w.Order = order
		artWorks = append(artWorks, w)
	}

	return artWorks
}

func (i instance) DownloadInventoryAsCSV(fileID string, destination string) error {
	response, err := i.drive.Files.Export(fileID, "text/csv").Download()
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
