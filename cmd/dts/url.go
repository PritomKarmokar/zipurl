package dts

type ShortUrlRequest struct {
	Url           string `json:"url" validate:"required,http_url"`
	Expiry        string `json:"expiry" validate:"omitempty,oneof=30m 60m 90m 7d 30d"`
	MaximumClicks *int64 `json:"maximum_clicks" validate:"omitempty,gte=1,lte=1000000"`
}
