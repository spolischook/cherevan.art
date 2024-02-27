package main

import (
	"github.com/cherevan.art/src/googleDrive/inventory"
	"github.com/kpango/glg"
)

func main() {
	s := inventory.New()
	//err := s.SetId(2, 888)
	//checkErr(err)
	artWorks := s.GetArtWorks()
	glg.Infof("Art works: %d", len(artWorks))
}

func checkErr(err error) {
	if err != nil {
		glg.Fatal(err)
	}
}
