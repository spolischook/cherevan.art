package googleDrive

import (
	"fmt"
	"github.com/cherevan.art/src/artWork"
	"github.com/cherevan.art/src/tool"
	"github.com/kpango/glg"
	"google.golang.org/api/drive/v3"
	"strings"
)

const paintingsFolderId = "1NA3FqpicdYnl7RS7-YDks0oh6eLNFTcp"
const graphicsFolderId = "11nADUmGwSlnT4LGA9DF7RzButDb1FHNd"
const berehyniFolderId = "1MAyXihXxeFFzXrIJSiFbydP1XBc1BJe6"
const womenDefendersFolderId = "1-kCaziyY87N-hLY831eg7j4MFA8NB2Yn"

var folderIds = []string{
	paintingsFolderId,
	graphicsFolderId,
	berehyniFolderId,
	womenDefendersFolderId}

func (s Service) ListAllOriginals() ([]*drive.File, error) {
	var originals []*drive.File
	for _, folderId := range folderIds {
		files, err := s.ListFilesInFolder(folderId)
		if err != nil {
			return nil, err
		}
		originals = append(originals, files...)
	}
	return originals, nil
}

func (s Service) ListFilesInFolder(folderId string) ([]*drive.File, error) {
	var files []*drive.File
	var pageToken string

	for {
		call := s.Files.List().Q(fmt.Sprintf("'%s' in parents", folderId))
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}

		r, err := call.Do()
		if err != nil {
			return nil, err
		}

		for _, i := range r.Files {
			files = append(files, i)
		}

		pageToken = r.NextPageToken
		if pageToken == "" {
			break
		}
	}

	return files, nil
}

func (s Service) FetchMainImage(aw artWork.ArtWork) error {
	if aw.ImageName == "" {
		return fmt.Errorf("No image name for art work: %#v\n", aw)
	}

	// check if the image already exists
	_, err := tool.FindFileInDir(aw.PageLeafPath(), aw.ImageName)
	if err == nil {
		return nil
	}

	// First check materials to determine if it's a painting or graphics
	var dirId string
	isPainting := false
	isGraphics := false

	for _, material := range aw.Materials {
		material = strings.ToLower(material)
		if strings.Contains(material, "canvas") || strings.Contains(material, "oil") {
			isPainting = true
			break
		}
		if strings.Contains(material, "paper") || strings.Contains(material, "cardboard") {
			isGraphics = true
			break
		}
	}

	// If materials don't give us a clear answer, check categories
	if !isPainting && !isGraphics {
		for _, cat := range aw.Categories {
			cat = strings.ToLower(cat)
			if cat == "painting" {
				isPainting = true
				break
			}
			if cat == "graphics" || cat == "pastel" {
				isGraphics = true
				break
			}
		}
	}

	// Set the directory based on our findings
	if isPainting {
		dirId = paintingsFolderId
	} else if isGraphics {
		dirId = graphicsFolderId
	} else {
		// Default to graphics if we can't determine
		dirId = graphicsFolderId
	}

	err = s.DownloadFile(aw.ImageName, dirId, aw.PageLeafPath())
	if err != nil {
		if dirId == graphicsFolderId {
			err = s.DownloadFile(aw.ImageName, paintingsFolderId, aw.PageLeafPath())
		}
		if err != nil {
			err = s.DownloadFile(aw.ImageName, berehyniFolderId, aw.PageLeafPath())
		}
		if err != nil {
			err = s.DownloadFile(aw.ImageName, womenDefendersFolderId, aw.PageLeafPath())
		}
	}

	if err != nil {
		glg.Errorf("Failed to find artwork '%s' (%s). Searched in:", aw.Title, aw.ImageName)
		if isPainting {
			glg.Errorf("- Paintings folder (based on materials: %v)", aw.Materials)
		} else if isGraphics {
			glg.Errorf("- Graphics folder (based on materials: %v)", aw.Materials)
		} else {
			glg.Errorf("- Graphics folder (default, no material match found in: %v)", aw.Materials)
		}
		glg.Errorf("- Bereyhni folder")
		glg.Errorf("- Women Defenders folder")
		glg.Errorf("Categories: %v", aw.Categories)
	}

	return err
}
