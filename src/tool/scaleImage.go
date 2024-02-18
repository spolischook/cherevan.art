package tool

import (
	"github.com/kpango/glg"
	"github.com/nfnt/resize"
	"image/jpeg"
	"os"
)

func ScaleImage(src string, width uint) {
	// Open the image file
	file, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Decode the image
	img, err := jpeg.Decode(file)
	if err != nil {
		glg.Fatalf("Error decoding image %s: %s\n", src, err)
	}

	// Scale the image
	scaledImg := resize.Resize(width, 0, img, resize.Lanczos3)

	// Create the output file
	out, err := os.Create(src)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// Write the scaled image to the output file
	jpeg.Encode(out, scaledImg, nil)
}
