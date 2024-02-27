package artWork

import (
	"fmt"
	"net/url"
	"strings"
)

func (aw ArtWork) GetUrl() string {
	return fmt.Sprintf("art-works/%s", aw.Slug)
}

func (aw ArtWork) ImageUrl() string {
	segments := strings.Split(aw.HugoUrl, "/")
	for i, seg := range segments {
		segments[i] = url.PathEscape(seg)
	}
	escapedHugoUrl := strings.Join(segments, "/")
	return fmt.Sprintf(
		"https://www.cherevan.art/%s/%s",
		//aw.HugoUrl,
		escapedHugoUrl,
		url.PathEscape(aw.ImageName))
}
