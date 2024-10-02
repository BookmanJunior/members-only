package cloud

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Cloudinary struct {
	Cld *cloudinary.Cloudinary
	Ctx context.Context
}

func (c *Cloudinary) Open(conString string) error {
	cld, err := cloudinary.NewFromURL(conString)

	if err != nil {
		return err
	}

	cld.Config.URL.Secure = true
	ctx := context.Background()
	c.Cld = cld
	c.Ctx = ctx

	return nil
}

func (c Cloudinary) UploadFile(userId int, fileName string) (string, error) {
	fName := fmt.Sprintf("members-only/attachments/%v/%v", userId, fileName)
	resp, err := c.Cld.Upload.Upload(c.Ctx, filepath.Join("./attachments", filepath.Base(fileName)), uploader.UploadParams{
		PublicID:       fName,
		UniqueFilename: api.Bool(false),
		Overwrite:      api.Bool(false),
	})

	if err != nil {
		return "", err
	}

	fmt.Println("**** Uploaded an image****\nDelivery URL:", resp.SecureURL)
	return resp.SecureURL, nil
}
