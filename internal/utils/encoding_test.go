package utils

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEncodedData(t *testing.T) {
	t.Parallel()

	t.Run("passing_test_gzip", func(t *testing.T) {
		t.Parallel()
		encodedStr := "eJwrSS0uAQAEXQHB"
		encodedBytes := []byte(encodedStr)
		reader := io.NopCloser(bytes.NewReader(encodedBytes))
		decodedBytes, err := GetEncodedData(reader, "gzip")
		assert.Nil(t, err)
		expected := []byte("test")
		assert.Equal(t, expected, decodedBytes)
	})
}
