package utils

import (
	"encoding/base64"
	uuid "github.com/satori/go.uuid"
	"strings"
)

var escaper = strings.NewReplacer("9", "99", "-", "90", "_", "91")
var unescaper = strings.NewReplacer("99", "9", "90", "-", "91", "_")

func ShortUuid(uuid uuid.UUID) string {
	return escaper.Replace(base64.RawURLEncoding.EncodeToString(uuid.Bytes()))
}

func ShortUuidFromStr(s string) (string, error) {
	uid, err := uuid.FromString(s)
	if err != nil {
		return "", err
	}
	return escaper.Replace(base64.RawURLEncoding.EncodeToString(uid.Bytes())), nil
}

func ParseShortUuid(str string) (uuid.UUID, error) {
	dec, err := base64.RawURLEncoding.DecodeString(unescaper.Replace(str))
	if err != nil {
		return uuid.Nil, err
	}
	uid, err := uuid.FromBytes(dec)
	if err != nil {
		return uuid.Nil, err
	}
	return uid, nil
}
