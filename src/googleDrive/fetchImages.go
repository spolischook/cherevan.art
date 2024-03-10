package googleDrive

import (
	"fmt"
	"github.com/cherevan.art/src/artWork"
	"github.com/cherevan.art/src/tool"
	"google.golang.org/api/drive/v3"
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

	// if ArtWork.Categories contains "painting" then use paintingsFolderId
	// else use graphicsFolderId
	var dirId string
	for _, cat := range aw.Categories {
		if cat == "painting" {
			dirId = paintingsFolderId
			break
		}
		if cat == "graphics" {
			dirId = graphicsFolderId
			break
		}
		if cat == "pastel" {
			dirId = graphicsFolderId
			break
		}
	}

	err = s.DownloadFile(aw.ImageName, dirId, aw.PageLeafPath())
	if err != nil {
		err = s.DownloadFile(aw.ImageName, berehyniFolderId, aw.PageLeafPath())
	}
	if err != nil {
		err = s.DownloadFile(aw.ImageName, womenDefendersFolderId, aw.PageLeafPath())
	}

	return err
}
