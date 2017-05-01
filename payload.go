package hooks

import (
	"net"
)

type PayloadWrapper struct {
	Payload
}

func (p PayloadWrapper) Interface() interface{} {
	if p.Payload == nil {
		return nil
	}
	return p.Payload
}

func (p PayloadWrapper) Err() error {
	if err, ok := p.Payload.(error); ok {
		return err
	}

	return nil
}

func (p PayloadWrapper) Listener() net.Listener {
	if ln, ok := p.Payload.(net.Listener); ok {
		return ln
	}

	return nil
}
