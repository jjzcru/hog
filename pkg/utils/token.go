package utils

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

func GetToken() string {
	hasher := md5.New()
	_, _ = hasher.Write([]byte(time.Now().Format(time.RFC3339)))
	return hex.EncodeToString(hasher.Sum(nil))[0:22]
}
