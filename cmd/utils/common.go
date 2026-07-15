package utils

import (
	cryptoRand "crypto/rand"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/spf13/viper"
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

func ParseExpiryDuration(expiry string) time.Duration {
	switch expiry {
	case "30m":
		return 30 * time.Minute
	case "60m":
		return 60 * time.Minute
	case "90m":
		return 90 * time.Minute
	case "7d":
		return 7 * 24 * time.Hour
	case "30d":
		return 30 * 24 * time.Hour
	default:
		return 30 * time.Minute
	}
}
