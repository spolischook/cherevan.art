package googleDrive

import (
	"context"
	"github.com/cherevan.art/src/googleDrive/login"
	"github.com/kpango/glg"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func New() (*Service, error) {
	tokenSource, err := login.TokenSource()
	if err != nil {
		glg.Fatal(err)
	}

	d, err := drive.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		glg.Fatal(err)
	}

	return &Service{d}, nil
}
