package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// generate test for the encoding function

func TestGetEncodedData(t *testing.T) {
	t.Parallel()

	t.Run("passing_test", func(t *testing.T) {
		t.Parallel()
		encodedStr := "eJwrSS0uAQAEXQHB"
		encodedBytes := []byte(encodedStr)
		decodedBytes, err := GetEncodedData(encodedBytes, "deflate")
		assert.Nil(t, err)
		expected := []byte("test")
		assert.Equal(t, expected, decodedBytes)
	})
}
