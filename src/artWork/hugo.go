package artWork

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func (a ArtWork) OldPath() string {
	return fmt.Sprintf("content/art-works/%s/%s", a.Date.Format("2006"), a.Slug)
}
func (a ArtWork) NewPath() string {
	return fmt.Sprintf("content/art-works/%s/%s", a.Date.Format("2006"), a.Slug)
}
func (a ArtWork) PageLeafPath() string {
	return a.NewPath()
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

func GetExistingArtWorks() (map[string]ArtWork, error) {
	artWorks := map[string]ArtWork{}

	err := filepath.Walk("content/art-works", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.Name() == "index.md" {
			artWork, err := NewArtWorkFromPath(path)
			if err != nil {
				return err
			}
			artWorks[path] = artWork
		}

		return nil
	})

	return artWorks, err
}

func NewArtWorkFromPath(path string) (ArtWork, error) {
	var artWork ArtWork

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return artWork, err
	}
	defer file.Close()

	// Read the file content
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return artWork, err
	}

	// Convert the content to a string
	contentStr := string(content)

	// Extract the YAML front matter
	start := strings.Index(contentStr, "---")
	end := strings.LastIndex(contentStr, "---")
	yamlContent := content[start+3 : end]

	// Unmarshal the YAML front matter into an ArtWork struct
	err = yaml.Unmarshal([]byte(yamlContent), &artWork)
	if err != nil {
		return artWork, err
	}

	// Extract the remaining content after the front matter
	textContent := strings.TrimSpace(contentStr[end+3:])

	// Assign the remaining content to the Text property of the ArtWork struct
	artWork.Text = textContent

	return artWork, nil
}
