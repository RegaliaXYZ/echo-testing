package utils

import (
	"context"
	"os"

	"cloud.google.com/go/storage"
)

type GoogleService struct {
	root *storage.BucketHandle
}

type GoogleServiceInterface interface {
	WriteFolder(name string) error
	DeleteFile(path string) error
	FileExists(path string) (bool, error)
	WriteFile(path string, v any) error
}

func NewGoogleService() GoogleServiceInterface {
	gcp_client, err := storage.NewClient(context.Background())
	if err != nil {
		panic("cannot create gcp client")
	}
	root_bucket := os.Getenv("ROOT_BUCKET")
	if root_bucket == "" {
		panic("ROOT_BUCKET not set")
	}
	gcp_handler := gcp_client.Bucket(root_bucket)

	return &GoogleService{
		root: gcp_handler,
	}
}

func (g *GoogleService) WriteFolder(name string) error {
	ctx := context.Background()
	folder := g.root.Object(name)
	if _, err := folder.Attrs(ctx); err != nil {
		if err == storage.ErrObjectNotExist {
			return folder.NewWriter(ctx).Close()
		} else {
			return err
		}
	}
	return nil
}

func (g *GoogleService) DeleteFile(path string) error {
	ctx := context.Background()
	return g.root.Object(path).Delete(ctx)
}

func (g *GoogleService) FileExists(path string) (bool, error) {
	ctx := context.Background()
	_, err := g.root.Object(path).Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (g *GoogleService) WriteFile(path string, v any) error {
	ctx := context.Background()

	w := g.root.Object(path).NewWriter(ctx)
	defer w.Close()
	if _, err := w.Write(v.([]byte)); err != nil {
		return err
	}
	return nil
}
