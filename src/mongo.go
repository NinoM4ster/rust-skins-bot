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

func upsertSkin(skin Skin) error {
	skins := mongoClient.Database("rust-skins").Collection("skins")
	marsh, err := bson.Marshal(skin)
	if err != nil {
		return err
	}
	ctx, cancel := newCtx(5)
	defer cancel()
	_, err = skins.ReplaceOne(ctx, bson.M{"page_url": skin.PageURL}, marsh, options.Replace().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}
