package util

import (
	"crypto/md5"
	"encoding/hex"
)

func GenSaltedMD5(input string, salt string) string {
	data := []byte(input + salt)
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}
