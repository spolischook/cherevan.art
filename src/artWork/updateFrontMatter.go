package artWork

func (aw *ArtWork) UpdateFrontMatter(a2 ArtWork) {
	aw.ID = a2.ID
	aw.Order = a2.Order
	aw.Title = a2.Title
	if aw.Slug != a2.Slug {
		aw.Aliases = append(aw.Aliases, aw.GetUrl())
	}
	aw.Slug = a2.Slug
	if aw.HugoUrl == "" {
		aw.HugoUrl = a2.GetUrl()
	}
	aw.Categories = a2.Categories
	aw.InStock = a2.InStock
	aw.IsVisible = a2.IsVisible
	aw.Height = a2.Height
	aw.Width = a2.Width
	aw.Date = a2.Date
	aw.Materials = a2.Materials
	aw.Price = a2.Price
	aw.ImageName = a2.ImageName
}
