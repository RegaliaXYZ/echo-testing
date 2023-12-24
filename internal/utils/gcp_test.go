package utils

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func TestGoogleService(t *testing.T) {
	t.Parallel()

	gcp_service := NewGoogleService()
	random_string := generateRandomString(10)
	t.Run("passing_test_write_folder", func(t *testing.T) {
		t.Parallel()
		err := gcp_service.WriteFolder(random_string)
		assert.Nil(t, err)
	})

	t.Run("passing_write_file", func(t *testing.T) {
		t.Parallel()
		err := gcp_service.WriteFile(random_string+"_test.txt", []byte("test"))
		assert.Nil(t, err)
	})

	t.Run("passing_test_file_exists", func(t *testing.T) {
		t.Parallel()
		exists, err := gcp_service.FileExists(random_string + "_test.txt")
		assert.Nil(t, err)
		assert.True(t, exists)
	})

	t.Run("passing_test_delete_file", func(t *testing.T) {
		t.Parallel()
		err := gcp_service.DeleteFile(random_string + "_test.txt")
		assert.Nil(t, err)
	})

	t.Run("passing_test_delete_folder", func(t *testing.T) {
		t.Parallel()
		err := gcp_service.DeleteFile(random_string)
		assert.Nil(t, err)
	})
}
