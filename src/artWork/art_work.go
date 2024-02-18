package artWork

import (
	"time"
)

type ArtWorks []ArtWork
type ArtWork struct {
	ID         int       `yaml:"id" json:"id"`
	Order      int       `yaml:"order" json:"order"`
	Title      string    `yaml:"title" json:"title"`
	Slug       string    `yaml:"slug" json:"slug"`
	Url        string    `yaml:"url" json:"url"`
	Aliases    []string  `yaml:"aliases" json:"aliases"`
	Categories []string  `yaml:"categories" json:"categories"`
	InStock    bool      `yaml:"inStock" json:"inStock"`
	IsVisible  bool      `yaml:"isVisible" json:"isVisible"`
	Height     int       `yaml:"height" json:"height"`
	Width      int       `yaml:"width" json:"width"`
	Date       time.Time `yaml:"date" json:"date"`
	Materials  []string  `yaml:"materials" json:"materials"`
	Price      int       `yaml:"price" json:"price"`
	ImageName  string    `yaml:"mainImage" json:"mainImage"`
	Text       string    `yaml:"-" json:"text"`
}
