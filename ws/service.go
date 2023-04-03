package ws

import (
	"go.uber.org/zap"
)

type WebSocketService struct {
	ClientList map[string]*WebSocketClient
	Logger     *zap.SugaredLogger
}
