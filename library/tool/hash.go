package tool

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"math/rand"
)

func Md5(s string) string {
	m := md5.New()
	m.Write([]byte(s))
	return hex.EncodeToString(m.Sum(nil))
}

func Sha256(str, salt string) string {
	h := sha1.New()
	h.Write([]byte(str))
	h.Write([]byte(salt))
	return hex.EncodeToString(h.Sum(nil))
}

func RandString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	_, _ = rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
