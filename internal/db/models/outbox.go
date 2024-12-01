package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Outbox struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	MessageID    string             `json:"messageId" bson:"message_id"`
	Topic        string             `json:"topic" bson:"topic"`
	Content      interface{}        `json:"content" bson:"content"`
	Status       string             `json:"status" bson:"status"`
	RetryCount   int                `json:"retryCount" bson:"retry_count"`
	CreatedAt    string             `json:"createdAt" bson:"created_at"`
	UpdatedAt    string             `json:"updatedAt" bson:"updated_at"`
	ErrorMessage string             `json:"errorMessage" bson:"error_message,omitempty"`
}
