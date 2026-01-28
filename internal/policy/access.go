package policy

import (
	"my-project/types"
	"github.com/google/uuid"
)

// ChatPolicy defines access control rules
type ChatPolicy struct{}

// CanReadMessage checks if user is a participant
// We attach it to ChatPolicy to match the usage in ChatLogic
func (p *ChatPolicy) CanReadMessage(userID uuid.UUID, msg types.Message) bool {
	return msg.SenderID == userID || msg.ReceiverID == userID
}