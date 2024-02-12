package login

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/kpango/glg"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

const SecretJson = "google_secret.json"
const TokenJson = "google_token.json"

func TokenSource() (oauth2.TokenSource, error) {
	b, err := getSecret()
	if err != nil {
		glg.Errorf("Unable to read client secret file: %v", err)
		return nil, err
	}

	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		glg.Errorf("Unable to parse client secret file to config: %v", err)
		return nil, err
	}

	tok, err := tokenFromFile(TokenJson)
	if err != nil {
		glg.Warnf("Token not found: %v", err)
		tok = getTokenFromWeb(config)
	}

	// Refresh the token if it's expired
	src := config.TokenSource(context.Background(), tok)
	tok, err = src.Token()
	if err != nil {
		tok = getTokenFromWeb(config)
	}

	saveToken(TokenJson, tok)

	return src, nil
}

func getSecret() ([]byte, error) {
	b, err := os.ReadFile(SecretJson)
	if err != nil {
		glg.Errorf("Unable to read client secret file: %v", err)
		return nil, err
	}
	return b, nil
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
