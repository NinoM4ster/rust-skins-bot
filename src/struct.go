package main

import "go.mongodb.org/mongo-driver/bson/primitive"

// Item .
type Item struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	ItemName    string             `bson:"item_name,omitempty"`
	DisplayName string             `bson:"display_name,omitempty"`
	PageURL     string             `bson:"page_url,omitempty"`
}

// Skin .
type Skin struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Num         int                `bson:"num,omitempty"`
	WorkshopID  string             `bson:"workshop_id,omitempty"`
	DisplayName string             `bson:"display_name,omitempty"`
	PageURL     string             `bson:"page_url,omitempty"`
	ImageURL    string             `bson:"image_url,omitempty"`
	ItemName    string             `bson:"item_name,omitempty"`
}
