package artWork

import (
	"gopkg.in/yaml.v2"
	"os"
)

func (aw ArtWork) Save() error {
	// Create the directory if it does not exist
	err := aw.createLeafFolder()
	if err != nil {
		return err
	}

	// Create the file
	file, err := os.Create(aw.NewPath() + "/index.md")
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the YAML front matter
	err = writeFrontMatter(aw, file)
	if err != nil {
		return err
	}

	// Write the text content
	_, err = file.Write([]byte(aw.Text))
	if err != nil {
		return err
	}

	return nil
}

func writeFrontMatter(a ArtWork, file *os.File) error {
	yamlData, err := yaml.Marshal(a)
	if err != nil {
		return err
	}
	_, err = file.Write([]byte("---\n"))
	if err != nil {
		return err
	}
	_, err = file.Write(yamlData)
	if err != nil {
		return err
	}
	_, err = file.Write([]byte("---\n"))
	if err != nil {
		return err
	}
	return nil
}

func (aw ArtWork) createLeafFolder() error {
	return os.MkdirAll(aw.NewPath(), 0755)
}
