package inbound

type ProcessMessagePort interface {
	ProcessMessage(msg string, source string) (string, error)
}
