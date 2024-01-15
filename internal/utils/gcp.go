package utils

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/storage"
)

type GoogleService struct {
	root *storage.BucketHandle
}

type GoogleServiceInterface interface {
	WriteFolder(path string) error
	DeleteFolder(path string) error
	WriteFile(path string, v any) error
	GetWriter(path string) *storage.Writer
	FileExists(path string) (bool, error)
	ReadFile(path string, v any) error
	DeleteFile(path string) error
}

func NewGoogleService(root_bucket string) GoogleServiceInterface {
	gcp_client, err := storage.NewClient(context.Background())
	if err != nil {
		panic("cannot create gcp client")
	}
	gcp_handler := gcp_client.Bucket(root_bucket)

	return &GoogleService{
		root: gcp_handler,
	}
}
func (g *GoogleService) GetWriter(path string) *storage.Writer {
	ctx := context.Background()
	return g.root.Object(path).NewWriter(ctx)
}

func (g *GoogleService) WriteFolder(path string) error {
	ctx := context.Background()
	folder := g.root.Object(path + "/")
	if _, err := folder.Attrs(ctx); err != nil {
		if err == storage.ErrObjectNotExist {
			return folder.NewWriter(ctx).Close()
		} else {
			return err
		}
	}
	return nil
}

func (g *GoogleService) DeleteFolder(path string) error {
	ctx := context.Background()
	return g.root.Object(path + "/").Delete(ctx)
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

func (g *GoogleService) ReadFile(path string, v any) error {
	ctx := context.Background()
	r, err := g.root.Object(path).NewReader(ctx)
	if err != nil {
		return err
	}
	defer r.Close()
	if err := json.NewDecoder(r).Decode(&v); err != nil {
		return err
	}
	return nil
}

func (g *GoogleService) WriteFile(path string, v any) error {
	ctx := context.Background()

	w := g.root.Object(path).NewWriter(ctx)
	defer w.Close()
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return err
	}
	return nil
}
