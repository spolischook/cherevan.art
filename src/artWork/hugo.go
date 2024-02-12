package artWork

import (
	"fmt"
	"os"
	"path/filepath"
)

func (a ArtWork) OldPath() string {
	return fmt.Sprintf("content/art-works/%s/%s", a.Date.Format("2006"), a.Slug)
}
func (a ArtWork) NewPath() string {
	return fmt.Sprintf("content/art-works/%s/%d-%s", a.Date.Format("2006"), a.ID, a.Slug)
}

func (a ArtWork) MoveToNewPath() error {
	_, err := os.Stat(a.NewPath())
	if err == nil {
		// already exists
		return nil
	}

	_, err = os.Stat(a.OldPath())
	if err != nil {
		// nothing to move
		return err
	}

	parentDir := filepath.Dir(a.NewPath())
	err = os.MkdirAll(parentDir, 0755)
	if err != nil {
		return err
	}

	// Move the directory
	err = os.Rename(a.OldPath(), a.NewPath())
	if err != nil {
		return err
	}

	return nil
}
