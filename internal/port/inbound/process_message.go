package inbound

type ProcessMessagePort interface {
	ProcessMessage(msg string) (string, error)
}
