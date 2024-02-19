package main

import (
	"github.com/cherevan.art/src/artWork"
	"github.com/cherevan.art/src/tool"
	"github.com/kpango/glg"
	"slices"
)

func checkErr(err error) {
	if err != nil {
		glg.Fatal(err)
	}
}

func main() {
	existingArtWorks, err := artWork.GetExistingArtWorks()
	checkErr(err)

	ids := []int{}
	for p, w := range existingArtWorks {
		if w.ID < 1 {
			glg.Errorf("ID is empty for ArtWork: %s", p)
		}
		if slices.Contains(ids, w.ID) {
			glg.Errorf("ID is not unique for ArtWork: %s", p)
		}
		ids = append(ids, w.ID)

		// check that w.ImageName is not empty and exists in the filesystem
		// if not, print the error and exit
		if w.ImageName == "" {
			glg.Errorf("ImageName is empty for ID: %d and Title: %s", w.ID, w.Title)
		}

		_, err := tool.FindFileInDir(w.PageLeafPath(), w.ImageName)
		if err != nil {
			glg.Errorf("Image file not found for ID: %d and Title: %s", w.ID, w.Title)
		}
	}
}
