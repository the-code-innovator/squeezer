package redis

import (
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"

	"github.com/the-code-innovator/squeezer/shortener"
)

type redisRepository struct {
	client *redis.Client
}

func NewRedisClient(redisURL string) (*redis.Client, error) {
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(options)
	_, err = client.Ping().Result()
	return client, err
}

func NewRedisRepository(redisURL string) (shortener.ShortLinkRepository, error) {
	repository := &redisRepository{}
	client, err := NewRedisClient(redisURL)
	if err != nil {
		return nil, errors.Wrap(err, "repository.redisRepository.newRedisRepository")
	}
	repository.client = client
	return repository, nil
}

func (r *redisRepository) generateKey(code string) string {
	return fmt.Sprintf("shortlink: %s", code)
}

func (r *redisRepository) Find(code string) (*shortener.ShortLink, error) {
	shortLink := &shortener.ShortLink{}
	key := r.generateKey(code)
	data, err := r.client.HGetAll(key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "repository.ShortLink.redisRepository.Find")
	}
	if len(data) == 0 {
		return nil, errors.Wrap(shortener.ErrorShortLinkNotFound, "repository.ShortLink.redisRepository.Find")
	}
	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "repository.ShortLink.redisRepository.Find")
	}
	shortLink.Code = data["code"]
	shortLink.URL = data["url"]
	shortLink.CreatedAt = createdAt
	return shortLink, nil
}

func (r *redisRepository) Store(shortLink *shortener.ShortLink) error {
	key := r.generateKey(shortLink.Code)
	data := map[string]interface{}{
		"code":       shortLink.Code,
		"url":        shortLink.URL,
		"created_at": shortLink.CreatedAt,
	}
	_, err := r.client.HMSet(key, data).Result()
	if err != nil {
		return errors.Wrap(err, "repository.ShortLink.redisRepository.Store")
	}
	return nil
}
