package utils

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"errors"
	"io"
)

var (
	ErrorEncodingNotSupported = errors.New("encoding not supported")
	InternalError             = errors.New("internal error")
)

func GetEncodedData(rawData []byte, header string) ([]byte, error) {
	var reader io.ReadCloser
	var err error
	// if content is sent from QC, should always be deflate since the encoding has been specified to not be default compressor
	switch header {
	case "gzip":

		// b64z := "eNqr5uVSAAKl5PzcgpzUklQlKwWlytRiJR2oeEFiZU5+YgpQOLpaqbSkBCRfmqGQnFiUolCUmpdanpijpKOglJlXkpoHlnT29/V19HNxjXd2AknkleYCRQ2BrJLKArDpJanFJZl56Uq1OgpwE/MLUvOAxqWV5qWgmhbkGuwaFOYa7+4aHOLp70fQxNpYXq5aAEtlPaI="
		z, _ := base64.StdEncoding.DecodeString(string(rawData))
		reader, err = zlib.NewReader(bytes.NewReader(z))
		if err != nil {
			return nil, InternalError
		}
	case "deflate":
		reader, err = zlib.NewReader(bytes.NewReader(rawData))
		if err != nil {
			return nil, InternalError
		}
	default:
		return nil, ErrorEncodingNotSupported
	}
	defer reader.Close()
	raw, err := io.ReadAll(reader)
	if err != nil {
		return nil, InternalError
	}
	return raw, nil
}
