package model

type GenerateShortUrlRequest struct {
	Url string `json:"url"`
}

type UpdateShortUrlRequest struct {
	ID   int64  `json:"id"`
	Url  string `json:"url"`
	Slug string `json:"slug"`
}
