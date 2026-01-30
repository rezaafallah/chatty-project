package consts

const (
	// Context Keys
	KeyUserID = "user_id"

	// Redis Channels
	TopicChatBroadcast = "chat.broadcast"
	QueueChatInbound   = "chat.inbound"

	// Error Messages
	ErrInternalServer     = "internal server error"
	ErrInvalidInput       = "invalid input parameters"
	ErrUnauthorized       = "unauthorized access"
	ErrEmptyContent       = "message content cannot be empty"
	ErrInvalidCredentials = "invalid credentials"
	ErrInvalidMnemonic    = "invalid mnemonic"
)