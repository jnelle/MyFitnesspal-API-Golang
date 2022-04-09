package utils

import b64 "encoding/base64"

func EncodeBase64(rawString string) string {
	return b64.StdEncoding.EncodeToString([]byte(rawString))
}

func DecodeBase64URL(rawString string) ([]byte, error) {
	return b64.RawURLEncoding.DecodeString(rawString)
}
