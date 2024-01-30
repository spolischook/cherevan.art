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

func (s *Service) DownloadSheetAsCSV(fileID string, destination string) error {
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

func (s *Service) DownloadFile(fileName, dirId, destDir string) error {
	r, err := s.Files.List().Q(
		fmt.Sprintf("name='%s' and '%s' in parents", fileName, dirId),
	).Do()
	if err != nil {
		return err
	}

	if len(r.Files) == 0 {
		return fmt.Errorf("no files found")
	}

	if len(r.Files) > 1 {
		glg.Fatalf("multiple files found for file: %s", fileName)
	}

	i := r.Files[0]
	resp, err := s.Files.Get(i.Id).Download()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath.Join(destDir, i.Name))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
