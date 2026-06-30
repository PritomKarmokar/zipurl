package dts

type ShortUrlRequest struct {
	Url string `json:"url" validate:"required"`
}
