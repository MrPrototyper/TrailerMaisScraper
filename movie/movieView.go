package movie

import (
	"context"
	"scrap-data/movie/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Save(client *mongo.Client, ctx context.Context, movie models.Movie) (primitive.ObjectID, error) {
	collection := client.Database("movies").Collection("tmp_movies")
	insertResult, err := collection.InsertOne(ctx, movie)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return insertResult.InsertedID.(primitive.ObjectID), nil
}
