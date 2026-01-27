package srv

import (
	"my-project/internal/core"
	// ... imports
)

type WorkerServer struct {
	Logic *core.Logic
}

func (w *WorkerServer) Start() {
	for {
		// BLPOP from Redis Queue
		// w.Logic.ProcessMessage(ctx, payload)
	}
}