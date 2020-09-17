package main

import "go.mongodb.org/mongo-driver/bson/primitive"

// Item .
type Item struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	ItemName    string             `bson:"item_name,omitempty"`
	DisplayName string             `bson:"display_name,omitempty"`
	PagePath    string             `bson:"page_path,omitempty"`
}

// Skin .
type Skin struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	ItemID      primitive.ObjectID `bson:"item_id,omitempty"`
	WorkshopID  int64              `bson:"workshop_id,omitempty"`
	DisplayName string             `bson:"display_name,omitempty"`
	PagePath    string             `bson:"page_path,omitempty"`
	ImageURL    string             `bson:"image_url,omitempty"`
}
