package main

import (
	"fmt"
	"github.com/cherevan.art/src/googleDrive"
	"io"
	"log"
	"os"
)

const fileName = "taming_the_dragon.JPG"
const paintingsFolderId = "1NA3FqpicdYnl7RS7-YDks0oh6eLNFTcp"
const graphicsFolderId = "11nADUmGwSlnT4LGA9DF7RzButDb1FHNd"

func main() {
	srv, err := googleDrive.CreateService()
	checkError(err)

	// Replace 'filename' with the name of the file you are looking for.
	r, err := srv.Files.List().Q(fmt.Sprintf("name='%s' and '%s' in parents", fileName, paintingsFolderId)).Do()
	//r, err := srv.Files.List().Q(fmt.Sprintf("name='%s'", fileName)).Do()
	checkError(err)

	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		fmt.Println("Files:")
		for _, i := range r.Files {
			fmt.Printf("%s (%s)\n", i.Name, i.Id)

			// Download the file
			resp, err := srv.Files.Get(i.Id).Download()
			if err != nil {
				log.Fatalf("Unable to download file: %v", err)
			}
			defer resp.Body.Close()

			// Save the file to disk
			out, err := os.Create(i.Name)
			if err != nil {
				log.Fatalf("Unable to create file: %v", err)
			}
			defer out.Close()

			_, err = io.Copy(out, resp.Body)
			if err != nil {
				log.Fatalf("Unable to write file: %v", err)
			}
		}
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
