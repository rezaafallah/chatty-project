package policy

import (
	"my-project/types"
	"github.com/google/uuid"
)

// CanReadMessage checks if user is a participant
func CanReadMessage(userID uuid.UUID, msg types.Message) bool {
	return msg.SenderID == userID || msg.ReceiverID == userID
}