package main

import (
	"github.com/cherevan.art/src/artWork"
	"github.com/cherevan.art/src/googleDrive"
	"github.com/cherevan.art/src/googleDrive/inventory"
	"github.com/cherevan.art/src/tool"
	"github.com/kpango/glg"
	"path/filepath"
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
	invent := inventory.New()

	artWorks := invent.GetArtWorks()
	glg.Infof("Art works: %d", len(artWorks))
	maxId := artWorks.MaxId()

	validArtWorks, validationErrors := artWorks.Validate(func(w artWork.ArtWork) error {
		return w.Validate()
	})
	glg.Warnf("Validation errors: %d", len(validationErrors))
	//printValidationErrors(validationErrors)

	worksWithoutId := validArtWorks.Filter(func(w artWork.ArtWork) bool { return w.ID < 1 })
	worksWithoutId.ForEach(func(w artWork.ArtWork) {
		maxId = maxId + 1
		err := invent.SetId(w.Order, maxId)
		if err != nil {
			glg.Error(err)
		}
	})

	newArtworks := artWork.ArtWorks{}
	for _, w := range validArtWorks {
		found := false
		for p, existing := range existingArtWorks {
			if w.ID == existing.ID {
				existing.UpdateFrontMatter(w)
				checkErr(existing.Save())
				delete(existingArtWorks, p)
				found = true
			}
		}
		if !found {
			newArtworks = append(newArtworks, w)
			checkErr(gDrive.FetchMainImage(w))
			checkErr(w.Save())
			imagePath := filepath.Join(w.PageLeafPath(), w.ImageName)
			tool.ScaleImage(imagePath, 1200)
		}
	}

	glg.Warnf("Unused existing art works: %d", len(existingArtWorks))
	for p, w := range existingArtWorks {
		glg.Warnf(`"%s - %s"`, w.Title, p)
	}
	glg.Warnf("New art works: %d", len(newArtworks))

	newArtworks.ForEach(func(w artWork.ArtWork) {
		glg.Warnf(`New ArtWork: "%s"`, w.Title)
	})
}

func printValidationErrors(errs map[*artWork.ArtWork]error) {
	if len(errs) == 0 {
		return
	}
	for w, err := range errs {
		if w.Title != "" {
			glg.Errorf(`"%s": %s`, w.Title, err)
		} else {
			glg.Error(err)
		}
	}
}
