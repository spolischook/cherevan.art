package googleDrive

import (
	"fmt"
	"github.com/kpango/glg"
	"io"
	"os"
	"path/filepath"

	"google.golang.org/api/drive/v3"
)

type Service struct {
	*drive.Service
}

func (s Service) DownloadSheetAsCSV(fileID string, destination string) error {
	response, err := s.Files.Export(fileID, "text/csv").Download()
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) DownloadFile(fileName, dirId, destDir string) error {
	r, err := s.Files.List().Q(
		fmt.Sprintf("name='%s' and '%s' in parents", fileName, dirId),
	).Do()
	if err != nil {
		glg.Errorf("Google Drive API error while searching for file '%s': %v", fileName, err)
		return err
	}

	if len(r.Files) == 0 {
		glg.Warnf("File '%s' not found in directory ID: %s", fileName, dirId)
		return fmt.Errorf(`file "%s" not found`, fileName)
	}

	if len(r.Files) > 1 {
		glg.Fatalf("Multiple files found with name '%s' in directory ID: %s", fileName, dirId)
	}

	i := r.Files[0]
	resp, err := s.Files.Get(i.Id).Download()
	if err != nil {
		glg.Errorf("Failed to download file '%s': %v", fileName, err)
		return err
	}
	defer resp.Body.Close()

	// create directory if it not exists
	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		glg.Errorf("Failed to create directory '%s': %v", destDir, err)
		return err
	}

	destPath := filepath.Join(destDir, i.Name)
	out, err := os.Create(destPath)
	if err != nil {
		glg.Errorf("Failed to create output file '%s': %v", destPath, err)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		glg.Errorf("Failed to write file contents for '%s': %v", fileName, err)
		return err
	}

	return err
}
