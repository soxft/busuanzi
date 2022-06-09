package jwtutil

import (
	"busuanzi/config"
	"busuanzi/library/tool"
	"strings"
)

// not a standard JWT, only to prevent the fake data

func Generate(userIdentity string) string {
	sign := tool.Sha256(userIdentity, config.Bsz.JwtSecret)
	return userIdentity + "." + sign
}

func Check(token string) string {
	arr := strings.Split(token, ".")
	if len(arr) != 2 {
		return ""
	}
	if sign := tool.Sha256(arr[0], config.Bsz.JwtSecret); sign != arr[1] {
		return ""
	}
	return arr[0]
}
