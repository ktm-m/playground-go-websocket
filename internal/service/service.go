package service

import (
	"github.com/ktm-m/playground-go-websocket/internal/port/inbound"
	"github.com/ktm-m/playground-go-websocket/internal/service/processmessage"
)

type Service struct {
	ProcessMessageService inbound.ProcessMessagePort
}

func NewService() *Service {
	return &Service{
		ProcessMessageService: processmessage.NewService(),
	}
}
