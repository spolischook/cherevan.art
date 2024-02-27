package artWork

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"time"
)

func (aw ArtWork) Validate() error {
	return validation.ValidateStruct(&aw,
		validation.Field(&aw.Title, validation.Required),
		validation.Field(&aw.Height, validation.Required, validation.Min(5)),
		validation.Field(&aw.Width, validation.Required, validation.Min(5)),
		validation.Field(&aw.Date, validation.Required, validation.Min(
			time.Date(1975, 1, 25, 0, 0, 0, 0, time.UTC))),
		validation.Field(&aw.ImageName, validation.Required),
	)
}

func (sl ArtWorks) Validate(f func(ArtWork) error) (ArtWorks, map[*ArtWork]error) {
	var res ArtWorks
	errs := map[*ArtWork]error{}
	for i := range sl {
		if err := f(sl[i]); err == nil {
			res = append(res, sl[i])
		} else {
			errs[&sl[i]] = err
		}
	}
	return res, errs
}
