package shortener

import (
	"errors"
	"time"

	errs "github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
)

var (
	ErrorShortLinkNotFound = errors.New("ShortLink: NOT FOUND")
	ErrorShortLinkInvalid  = errors.New("ShortLink: INVALID")
)

type shortLinkService struct {
	shortLinkRepository ShortLinkRepository
}

func NewShortLinkService(shortLinkRepository ShortLinkRepository) ShortLinkService {
	return &shortLinkService{
		shortLinkRepository: shortLinkRepository,
	}
}

func (s *shortLinkService) Find(code string) (*ShortLink, error) {
	return s.shortLinkRepository.Find(code)
}

func (s *shortLinkService) Store(shortLink *ShortLink) error {
	if err := validate.Validate(shortLink); err != nil {
		return errs.Wrap(ErrorShortLinkInvalid, "service.ShortLink.Store")
	}
	shortLink.Code = shortid.MustGenerate()
	shortLink.CreatedAt = time.Now().UTC().Unix()
	return s.shortLinkRepository.Store(shortLink)
}
