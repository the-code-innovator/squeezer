package shortener

type ShortLinkRepository interface {
	Find(code string) (*ShortLink, error)
	Store(shortlink *ShortLink) error
}
