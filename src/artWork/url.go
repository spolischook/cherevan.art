package artWork

import (
	"fmt"
)

func (a ArtWork) GetUrl() string {
	return fmt.Sprintf("art-works/%s", a.Slug)
}
