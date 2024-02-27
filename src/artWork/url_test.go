package artWork

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var TestImageUrlTC = map[string]ArtWork {
	"https://www.cherevan.art/art-works/almost-naked/Almost%20naked.jpg": {
		HugoUrl:   "art-works/almost-naked",
		ImageName: "Almost naked.jpg",
	},
	"https://www.cherevan.art/art-works/%D1%81ontract-with-the-devil/contract_with_the_devil.jpg": {
		HugoUrl:   "art-works/сontract-with-the-devil",
		ImageName: "contract_with_the_devil.jpg",
	},
}

func TestImageUrl(t *testing.T) {
	for expected, aw := range TestImageUrlTC {
		assert.Equal(t, expected, aw.ImageUrl(), "Expected %s, got %s", expected, aw.ImageUrl())
	}
}
