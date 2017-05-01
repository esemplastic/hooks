package hooks

import (
	"net"
)

type SinglePayloadWrapper struct {
	Payload
}

func (p SinglePayloadWrapper) Interface() interface{} {
	return p.Payload
}

func (p SinglePayloadWrapper) Err() error {
	if err, ok := p.Payload.(error); ok {
		return err
	}

	return nil
}

func (p SinglePayloadWrapper) Listener() net.Listener {
	if ln, ok := p.Payload.(net.Listener); ok {
		return ln
	}

	return nil
}

type PayloadWrapper struct {
	payloads []Payload
}

func (p PayloadWrapper) Index(idx int) SinglePayloadWrapper {
	if idx+1 > len(p.payloads) {
		return SinglePayloadWrapper{}
	}

	return SinglePayloadWrapper{Payload: p.payloads[idx]}
}

// including start, excluding end
func (p PayloadWrapper) Range(start int, end int) PayloadWrapper {
	if end >= len(p.payloads) {
		return PayloadWrapper{}

	}
	return PayloadWrapper{payloads: []Payload{p.payloads[start:end]}}
}

func (p PayloadWrapper) First() SinglePayloadWrapper {
	return p.Index(0)
}

func (p PayloadWrapper) Second() SinglePayloadWrapper {
	return p.Index(1)
}

func (p PayloadWrapper) Last() SinglePayloadWrapper {
	return p.Index(len(p.payloads) - 1)
}
