package shortener

type ShortLinkSerializer interface {
	Decode(input []byte) (*ShortLink, error)
	Encode(input *ShortLink) ([]byte, error)
}
