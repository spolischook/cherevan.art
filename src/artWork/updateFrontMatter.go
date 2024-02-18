package artWork

func (a *ArtWork) UpdateFrontMatter(a2 ArtWork) {
	a.ID = a2.ID
	a.Order = a2.Order
	a.Title = a2.Title
	if a.Slug != a2.Slug {
		a.Aliases = append(a.Aliases, a.GetUrl())
	}
	a.Slug = a2.Slug
	if a.Url == "" {
		a.Url = a2.GetUrl()
	}
	a.Categories = a2.Categories
	a.InStock = a2.InStock
	a.IsVisible = a2.IsVisible
	a.Height = a2.Height
	a.Width = a2.Width
	a.Date = a2.Date
	a.Materials = a2.Materials
	a.Price = a2.Price
	a.ImageName = a2.ImageName
}
