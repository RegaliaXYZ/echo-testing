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
	ErrorInternal             = errors.New("internal error")
)

func GetEncodedData(body io.ReadCloser, header string) ([]byte, error) {
	var reader io.ReadCloser
	var err error
	rawData, err := io.ReadAll(body)
	if err != nil {
		return nil, ErrorInternal
	}
	// if content is sent from QC, should always be deflate since the encoding has been specified to not be default compressor
	if header != "gzip" {
		return nil, ErrorEncodingNotSupported
	}
	// b64z := "eNqr5uVSAAKl5PzcgpzUklQlKwWlytRiJR2oeEFiZU5+YgpQOLpaqbSkBCRfmqGQnFiUolCUmpdanpijpKOglJlXkpoHlnT29/V19HNxjXd2AknkleYCRQ2BrJLKArDpJanFJZl56Uq1OgpwE/MLUvOAxqWV5qWgmhbkGuwaFOYa7+4aHOLp70fQxNpYXq5aAEtlPaI="
	z, _ := base64.StdEncoding.DecodeString(string(rawData))
	reader, err = zlib.NewReader(bytes.NewReader(z))
	if err != nil {
		return nil, ErrorInternal
	}

	defer reader.Close()
	raw, err := io.ReadAll(reader)
	if err != nil {
		return nil, ErrorInternal
	}
	return raw, nil
}
