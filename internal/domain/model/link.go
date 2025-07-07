package model

type Link struct {
	ID          int
	OriginalUrl string
	ShortUrl    string
}

func New(url string) *Link {
	return &Link{ShortUrl: url}
}
