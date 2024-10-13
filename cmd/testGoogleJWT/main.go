package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Route to handle POST request to receive JWT token
	router.POST("/token", func(c *gin.Context) {
		var json struct {
			Token string `json:"token" binding:"required"`
		}

		// Bind incoming JSON data
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Log the received token
		fmt.Println("Received JWT token:", json.Token)

		// Use the token to access the Google Sheet and convert to CSV
		csvData, err := getSheetDataAsCSV(json.Token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = ioutil.WriteFile("/tmp/inventory.csv", []byte(csvData), 0644)
		if err != nil {
			panic(err)
		}

		fmt.Println("CSV data saved to /tmp/inventory.csv")

		// Return the CSV data as a response
		c.Header("Content-Type", "text/csv")
		c.String(http.StatusOK, "")
	})

	// Start the server
	router.Run(":8080")
}

// getSheetDataAsCSV retrieves Google Sheet data using an access token and converts it to CSV format
func getSheetDataAsCSV(token string) (string, error) {
	// Replace with your specific Google Sheets API endpoint for exporting the spreadsheet as CSV
	sheetID := "1MPkRsuhJhJdYlTCOxHFiJoFX3DZalZPbTOUObacJPg0"
	gid := "1746605435"
	url := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/export?format=csv&gid=%s", sheetID, gid)

	// Make an authenticated GET request to Google Sheets API using the provided access token
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// Set the Authorization header with the access token
	req.Header.Set("Authorization", "Bearer "+token)

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check if the response was successful
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to fetch Google Sheet: %s", string(bodyBytes))
	}

	// Read the response body (CSV data)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
