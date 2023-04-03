package main

import (
	"WsGolangServer/ws"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewSugaredLogger(logConfig *zap.Config) (*zap.SugaredLogger, error) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.OutputPaths = logConfig.OutputPaths
	zapLogger, err := cfg.Build()
	if err != nil {
		return nil, err
	} else {
		return zapLogger.Sugar(), nil
	}
}

type ServerConfig struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

func (v *ServerConfig) getServerAddress() string {
	return v.Address + ":" + strconv.Itoa(v.Port)
}

func main() {
	logConfig := zap.Config{
		OutputPaths: []string{"logs/ws.log"},
	}
	serverConfig := ServerConfig{
		Address: "127.0.0.1",
		Port:    8080,
	}

	logger, _ := NewSugaredLogger(&logConfig)

	wsServer := ws.NewWebSocketServer(serverConfig.getServerAddress())
	wsServerService := ws.WebSocketService{
		Logger: logger,
	}
	wsServerService.StartServer(*wsServer)

	wsServer.Listen()
}
