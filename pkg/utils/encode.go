package utils

import (
	"encoding/base64"
)

func Encode(s string) (string, error) {
	encoded := base64.StdEncoding.EncodeToString([]byte(s))
	return encoded, nil
}

func Decode(s string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
