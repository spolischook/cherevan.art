package sheets

import (
	"context"
	"encoding/json"
	"github.com/cherevan.art/src/googleDrive"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"os"
	"strconv"
)

type token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	Expiry       string `json:"expiry"`
}

type instance struct {
	*sheets.Service
	spreadsheetId string
}

func New(spreadsheetId string) instance {
	// Read the token from the file
	file, err := os.Open(googleDrive.TokenJson)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	tok := &token{}
	err = json.NewDecoder(file).Decode(tok)
	if err != nil {
		panic(err)
	}

	// Create a TokenSource from the token
	t := &oauth2.Token{
		AccessToken:  tok.AccessToken,
		TokenType:    tok.TokenType,
		RefreshToken: tok.RefreshToken,
	}
	tokenSource := oauth2.StaticTokenSource(t)

	// Use the TokenSource to authenticate the Sheets service
	s, err := sheets.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		panic(err)
	}

	return instance{s, spreadsheetId}
}

func (i instance) SetId(rowNumber, id int) error {
	vr := &sheets.ValueRange{
		Values: [][]interface{}{{id}},
	}

	cellRange := "Sheet1!A" + strconv.Itoa(rowNumber+1)

	_, err := i.Spreadsheets.Values.Update(i.spreadsheetId, cellRange, vr).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}
