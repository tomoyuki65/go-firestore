package message

import (
	"time"
)

type MessageDocument struct {
	SenderID  string    `firestore:"senderId"`
	Text      string    `firestore:"text"`
	CreatedAt time.Time `firestore:"createdAt"`
	UpdatedAt time.Time `firestore:"updatedAt"`
}
