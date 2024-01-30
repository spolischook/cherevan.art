package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	mysqlUsername  = "root"
	mysqlPassword  = "pass"
	mysqlHost      = "127.0.0.1:3406"
	mysqlDatabase  = "cherry_t_art_com"
	hugoContentDir = "/home/spolischook/www/cherevan.art/content/archive"
)

func main() {
	// Create MySQL connection
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s", mysqlUsername, mysqlPassword, mysqlHost, mysqlDatabase))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Query to select data from jos_content table
	query := "SELECT id, title, alias, introtext, `fulltext`, created FROM jos_content"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate through the rows and create Hugo Markdown files
	for rows.Next() {
		var id int
		var title, alias, introtext, fulltext string
		var created []byte

		err := rows.Scan(&id, &title, &alias, &introtext, &fulltext, &created)
		if err != nil {
			log.Fatal(err)
		}

		// replace duble quotes with single
		title = strings.Replace(title, "\"", "'", -1)

		// Create a directory for each content item based on its ID
		contentDir := filepath.Join(hugoContentDir, fmt.Sprintf("%d-%s", id, alias))
		os.MkdirAll(contentDir, os.ModePerm)

		// Used to parse MySQL DATETIME into Go time.Time
		createdString := string(created)
		createdTime, err := time.Parse("2006-01-02 15:04:05", createdString)
		if err != nil {
			log.Fatal(err)
		}

		// Create Hugo Markdown file
		filePath := filepath.Join(contentDir, "index.md")
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Write Hugo front matter
		file.WriteString("---\n")
		file.WriteString(fmt.Sprintf("title: \"%s\"\n", title))
		file.WriteString(fmt.Sprintf("date: %s\n", createdTime.Format("2006-01-02 15:04:05")))
		file.WriteString(fmt.Sprintf("slug: %s\n", alias))
		file.WriteString("draft: true\n")
		file.WriteString("---\n\n")

		// Write content
		file.WriteString("{{< rawhtml >}}\n")
		file.WriteString(introtext)
		file.WriteString(fulltext)
		file.WriteString("\n{{< /rawhtml >}}\n")

		fmt.Printf("Exported content %d to %s\n", id, filePath)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Export complete.")
}
