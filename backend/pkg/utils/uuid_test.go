package utils

import (
	uuid "github.com/satori/go.uuid"
	"testing"
)

func TestEncoding(t *testing.T) {
	uid := uuid.NewV4()
	sh := ShortUuid(uid)
	println(sh)
}

func TestDecoding(t *testing.T) {
	uid := uuid.NewV4()
	println(uid.String())
	sh := ShortUuid(uid)
	println(sh)
	uid, err := ParseShortUuid(sh)
	if err != nil {
		panic(err)
	}
	println(uid.String())
}
