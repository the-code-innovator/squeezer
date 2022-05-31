package mongodb

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/the-code-innovator/squeezer/shortener"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

func NewMongoClient(mongodbURL string, mongodbTimeOut int) (*mongo.Client, error) {
	context, cancel := context.WithTimeout(context.Background(), time.Duration(mongodbTimeOut)*time.Second)
	defer cancel()
	client, err := mongo.Connect(context, options.Client().ApplyURI(mongodbURL))
	if err != nil {
		return nil, err
	}
	err = client.Ping(context, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, err
}

func NewMongoRepository(mongodbURL, mongoDB string, mongodbTimeOut int) (shortener.ShortLinkRepository, error) {
	repository := &mongoRepository{
		timeout:  time.Duration(mongodbTimeOut) * time.Second,
		database: mongoDB,
	}
	client, err := NewMongoClient(mongodbURL, mongodbTimeOut)
	if err != nil {
		return nil, errors.Wrap(err, "repository.mongoRepository.newMongoRepository")
	}
	repository.client = client
	return repository, nil
}

func (m *mongoRepository) Find(code *string) (*shortener.ShortLink, error) {
	context, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	shortLink := &shortener.ShortLink{}
	collection := m.client.Database(m.database).Collection("shortlinks")
	filter := bson.M{"code": code}
	err := collection.FindOne(context, filter).Decode(&shortLink)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(shortener.ErrorShortLinkNotFound, "repository.ShortLink.mongoRepository.Find")
		}
		return nil, errors.Wrap(err, "repository.ShortLink.mongoRepository.Find")
	}
	return shortLink, nil
}

func (m *mongoRepository) Store(shortLink *shortener.ShortLink) error {
	context, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	collection := m.client.Database(m.database).Collection("shortlinks")
	_, err := collection.InsertOne(
		context,
		bson.M{
			"code":       shortLink.Code,
			"created_at": shortLink.CreatedAt,
			"url":        shortLink.URL,
		},
	)
	if err != nil {
		return errors.Wrap(err, "repository.ShortLink.mongoRepository.Store")
	}
	return nil
}
