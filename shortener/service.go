package shortener

type ShortLinkService interface {
	Find(code *string) (*ShortLink, error)
	Store(shortlink *ShortLink) error
}
