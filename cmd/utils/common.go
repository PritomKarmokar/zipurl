package utils

import (
	cryptoRand "crypto/rand"
	"github.com/oklog/ulid/v2"
	"time"
)

func GenerateULID() string {
	t := time.Now()
	entropy := ulid.Monotonic(cryptoRand.Reader, 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}
