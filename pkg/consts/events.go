package consts

const (
	// Client -> Server Events
	EventMessageNew    = "message.new"
	EventMessageEdit   = "message.edit"
	EventMessageDelete = "message.delete"
	EventUserTyping    = "user.typing"

	// Server -> Client Events
	EventMessageAck = "message.ack"
)