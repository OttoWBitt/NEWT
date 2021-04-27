package main

import (
	"encoding/base64"
)

func encodeToString(enc string) string {
	data := []byte(enc)
	str := base64.StdEncoding.EncodeToString(data)
	return str
}

func decodeString(dec string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(dec)
	if err != nil {
		return "", err
	}
	return bytesToString(data), nil
}

func bytesToString(data []byte) string {
	return string(data[:])
}
