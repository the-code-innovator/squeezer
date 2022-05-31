package msgpack

import (
	"github.com/pkg/errors"
	"github.com/the-code-innovator/squeezer/shortener"
	"github.com/vmihailenco/msgpack"
)

type ShortLink struct{}

func (s *ShortLink) Decode(input []byte) (*shortener.ShortLink, error) {
	shortLink := &shortener.ShortLink{}
	if err := msgpack.Unmarshal(input, shortLink); err != nil {
		return nil, errors.Wrap(err, "serializer.ShortLink.Decode")
	}
	return shortLink, nil
}

func (s *ShortLink) Encode(input *shortener.ShortLink) ([]byte, error) {
	rawMessage, err := msgpack.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.ShortLink.Encode")
	}
	return rawMessage, nil
}
