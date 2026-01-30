package network

import (
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// SocketConnection defines the interface we use across the app
// so we can swap implementations if needed.
type SocketConnection interface {
	WriteMessage(messageType int, data []byte) error
	ReadMessage() (messageType int, p []byte, err error)
	Close() error
	SetReadLimit(limit int64)
	SetReadDeadline(t time.Time) error
	SetPongHandler(h func(appData string) error)
	SetWriteDeadline(t time.Time) error
	// Returns a standard io.WriteCloser to make streaming writes easier
	NextWriter(messageType int) (io.WriteCloser, error)
}

// WebSocketConn wraps gorilla's connection
// and implements SocketConnection.
type WebSocketConn struct {
	*websocket.Conn
}

// Re-export common websocket constants so callers
// don’t need to import gorilla directly.
const (
	CloseMessage         = websocket.CloseMessage
	TextMessage          = websocket.TextMessage
	PingMessage          = websocket.PingMessage
	PongMessage          = websocket.PongMessage
	CloseGoingAway       = websocket.CloseGoingAway
	CloseAbnormalClosure = websocket.CloseAbnormalClosure
)

// Upgrader Wrapper
type Upgrader struct {
	internal websocket.Upgrader
}

func NewUpgrader(readBuf, writeBuf int) *Upgrader {
	return &Upgrader{
		internal: websocket.Upgrader{
			ReadBufferSize:  readBuf,
			WriteBufferSize: writeBuf,
			// Allow all origins for now.
			// Tighten this in production if needed.
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
}

func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*WebSocketConn, error) {
	conn, err := u.internal.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}
	return &WebSocketConn{Conn: conn}, nil
}

// Helper wrapper so callers don’t depend on gorilla directly.
func IsUnexpectedCloseError(err error, closeCodes ...int) bool {
	return websocket.IsUnexpectedCloseError(err, closeCodes...)
}

// NextWriter passes through to the underlying connection
// but returns a standard io.WriteCloser.
func (c *WebSocketConn) NextWriter(messageType int) (io.WriteCloser, error) {
	return c.Conn.NextWriter(messageType)
}