package artWork

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"time"
)

func (a ArtWork) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Title, validation.Required),
		validation.Field(&a.Height, validation.Required, validation.Min(5)),
		validation.Field(&a.Width, validation.Required, validation.Min(5)),
		validation.Field(&a.Date, validation.Required, validation.Min(
			time.Date(1975, 1, 25, 0, 0, 0, 0, time.UTC))),
		validation.Field(&a.ImageName, validation.Required),
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
