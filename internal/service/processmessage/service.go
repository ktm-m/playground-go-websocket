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

func (s *service) ProcessMessage(msg string, source string) (string, error) {
	log.Printf("[SERVICE] message received from %s: %s", source, msg)
	return msg, nil
}
