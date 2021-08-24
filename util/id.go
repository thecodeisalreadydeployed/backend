package util

import (
	"strings"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

func RandomString(length int) string {
	id, err := gonanoid.Generate("-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", length)
	if err != nil {
		panic(err)
	}
	return strings.ToLower(id)
}
