package googleDrive

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var (
	instance *sheets.Service
	once     sync.Once
)

func (s Service) getSheetsService() (*sheets.Service, error) {
	var err error
	once.Do(func() {
		instance, err = sheets.NewService(context.Background(), option.WithCredentialsFile("path_to_credentials.json"))
	})
	return instance, err
}

func (s Service) SetId(rowNumber, id int) error {
	// Create a ValueRange to hold the new ID
	vr := &sheets.ValueRange{
		Values: [][]interface{}{{id}},
	}

	// Convert the row number and column number (assuming column 1 for ID) to an A1 notation range
	// Note: rowNumber + 1 is used because rows are 1-indexed in A1 notation
	cellRange := fmt.Sprintf("Sheet1!A%d", rowNumber+1)

	// Get the singleton Sheets service
	srv, err := s.getSheetsService()
	if err != nil {
		return fmt.Errorf("error getting Sheets service: %v", err)
	}

	// Call the Sheets API's Update method to update the cell
	_, err = srv.Spreadsheets.Values.Update(inventoryFileID, cellRange, vr).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return fmt.Errorf("error updating cell: %v", err)
	}

	return nil
}
