package server

type SuccessResponse struct {
	ShortUrl string `json:"short_url"`
}
type ShortUrlResponse struct {
	ShortUrl string `json:"short_url"`
	Url      string `json:"url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
