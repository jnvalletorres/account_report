package repository

import (
	"account_report/adapter/repository/document"
	"account_report/domain/entity"
	"account_report/domain/port"
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoNotificationRepository struct {
	client  *mongo.Client
	db      string
	timeout time.Duration
	logger  port.LoggerPort
}

const COLLECTION string = "notifications"

func newMongClient(mongoServerURL string, timeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)

	defer cancel()

	client, err := mongo.Connect(
		ctx,
		options.Client().
			ApplyURI(mongoServerURL),
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewMongoNotificationRepository(mongoServerURL, mongoDb string, timeout int, logger port.LoggerPort) (port.NotificationRepositoryPort, error) {
	mongoClient, err := newMongClient(mongoServerURL, timeout)
	repo := &mongoNotificationRepository{
		client:  mongoClient,
		db:      mongoDb,
		timeout: time.Duration(timeout) * time.Second,
		logger:  logger,
	}
	if err != nil {
		return nil, errors.Wrap(err, "client error")
	}

	return repo, nil

}

func (r *mongoNotificationRepository) SaveNotification(notificationEntity entity.NotificationEntity) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	collection := r.client.Database(r.db).Collection(COLLECTION)
	notificationDocument := document.ToCreateNotificationDocument(notificationEntity)

	insertResult, err := collection.InsertOne(
		ctx,
		notificationDocument,
	)
	if err != nil {
		return "", err
	}

	id := insertResult.InsertedID.(primitive.ObjectID).Hex()

	return id, err
}

func (r *mongoNotificationRepository) FindById(id string) (*entity.NotificationEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)

	defer cancel()

	notificationDocument := &document.NotificationDocument{}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	collection := r.client.Database(r.db).Collection(COLLECTION)
	err = collection.FindOne(ctx, filter).Decode(notificationDocument)
	if err != nil {
		return nil, err
	}

	return document.ToNotificationEntity(*notificationDocument), nil
}

func (r *mongoNotificationRepository) UpdateById(id string, notificationEntity entity.NotificationEntity) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	collection := r.client.Database(r.db).Collection(COLLECTION)
	document := document.ToUpdateNotificationDocument(notificationEntity)
	result, err := collection.UpdateOne(ctx, filter, bson.M{"$set": &document})
	if err != nil {
		r.logger.Error("", err)
		return err
	}
	if result.MatchedCount == 0 {
		msg := fmt.Sprintf("Notification with filter: %s not found", filter)
		r.logger.Error(msg)
	}

	return nil
}
