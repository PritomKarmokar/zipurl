package utils

import (
	cryptoRand "crypto/rand"
	"github.com/oklog/ulid/v2"
	"github.com/spf13/viper"
	"strings"
	"time"
)

const (
	DisplayTimeFormat = "02 Jan 2006 03:04:05 PM"
)

func GenerateULID() string {
	t := time.Now()
	entropy := ulid.Monotonic(cryptoRand.Reader, 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

func FormatTime(t time.Time) string {
	timeZone, _ := time.LoadLocation(viper.GetString("TIME_ZONE"))
	return t.In(timeZone).Format(DisplayTimeFormat)
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
