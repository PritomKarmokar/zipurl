package utils

import (
	cryptoRand "crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/spf13/viper"
)

const (
	DisplayTimeFormat = "02 Jan 2006 03:04:05 PM"
	base62Chars       = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	tokenLength       = 6 // 6 chars gives ~56.8 billion combinations
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

var allowedExpiryDurations = map[string]time.Duration{
	"30m": 30 * time.Minute,
	"60m": 60 * time.Minute,
	"90m": 90 * time.Minute,
	"7d":  7 * 24 * time.Hour,
	"30d": 30 * 24 * time.Hour,
}

func ParseExpiryDuration(expiry string) (time.Duration, error) {
	duration, ok := allowedExpiryDurations[expiry]
	if !ok {
		return 0, fmt.Errorf("invalid expiry value: %q", expiry)
	}

	return duration, nil
}

func GenerateShortToken() string {
	result := make([]byte, tokenLength)
	maxVal := big.NewInt(int64(len(base62Chars)))

	for i := range result {
		n, err := cryptoRand.Int(cryptoRand.Reader, maxVal)
		if err != nil {
			panic(err)
		}
		result[i] = base62Chars[n.Int64()]
	}
	return string(result)
}
