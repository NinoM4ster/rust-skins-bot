package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func newCtx(t time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), t*time.Second)
}

// func insertSkin(skin Skin) error {
// 	ctx, cancel := newCtx(5)
// 	defer cancel()
// 	skins := mongoClient.Database("rust-skins").Collection("skins")
// 	_, err := skins.InsertOne(ctx, skin)
// 	return err
// }

func skinExists(skin Skin) bool {
	skins := mongoClient.Database("rust-skins").Collection("skins")
	ctx, cancel := newCtx(5)
	defer cancel()
	count, err := skins.CountDocuments(ctx, bson.M{"page_url": skin.PageURL})
	if err != nil {
		return false
	}
	if count == 0 {
		return false
	}
	return true
}

func upsertSkin(skin Skin) error {
	skins := mongoClient.Database("rust-skins").Collection("skins")
	ctx, cancel := newCtx(5)
	defer cancel()
	_, err = skins.ReplaceOne(ctx, bson.M{"page_url": skin.PageURL}, skin, options.Replace().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

func getLastNum() (int, error) {
	skins := mongoClient.Database("rust-skins").Collection("skins")
	var skin Skin
	ctx, cancel := newCtx(5)
	defer cancel()
	err := skins.FindOne(ctx, bson.M{}, options.FindOne().SetSort(bson.M{"num": -1})).Decode(&skin)
	return skin.Num, err
}
