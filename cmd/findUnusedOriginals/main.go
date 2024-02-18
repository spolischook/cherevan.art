package main

import (
	"fmt"
	"github.com/cherevan.art/src/artWork"
	"github.com/cherevan.art/src/googleDrive"
	"github.com/kpango/glg"
	"path/filepath"
	"slices"
	"strings"
)

func checkErr(err error) {
	if err != nil {
		glg.Fatal(err)
	}
}

func main() {
	existingArtWorks, err := artWork.GetExistingArtWorks()
	checkErr(err)

	gDrive, err := googleDrive.New()
	checkErr(err)
	allOriginals, err := gDrive.ListAllOriginals()
	checkErr(err)

	glg.Infof("All originals: %d", len(allOriginals))

	skipExtensions := []string{".pdf", ".nef", ".txt", ".tiff", ".heic", ""}
	for _, o := range allOriginals {
		if slices.Contains(skipExtensions, strings.ToLower(filepath.Ext(o.Name))) {
			continue
		}

		found := false
		for p, existing := range existingArtWorks {
			if o.Name == existing.ImageName {
				delete(existingArtWorks, p)
				found = true
				//glg.Infof("Used original: %s", o.Name)
			}
		}
		if !found {
			fmt.Println(o.Name)
		}
	}

}
