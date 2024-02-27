package artWork

import validation "github.com/go-ozzo/ozzo-validation"

func (aw ArtWork) Salable() error {
	rules := []*validation.FieldRules{
		validation.Field(&aw.InStock, validation.Required),
		validation.Field(&aw.Price, validation.Min(1)),
	}

	return validation.ValidateStruct(&aw, rules...)
}
