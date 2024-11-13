package processmessage

import (
	"github.com/ktm-m/playground-go-websocket/internal/port/inbound"
	"log"
)

type service struct {
}

func NewService() inbound.ProcessMessagePort {
	return &service{}
}

func (s *service) ProcessMessage(msg string) (string, error) {
	log.Println("[SERVICE] message received:", msg)
	return msg, nil
}
