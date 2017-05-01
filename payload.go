package hooks

import (
	"fmt"
	"net"
	"reflect"
	"strconv"
)

type SinglePayloadWrapper struct {
	Payload
	typ reflect.Type
}

func (p SinglePayloadWrapper) IsNil() bool {
	return p.Payload == nil
}

func (p SinglePayloadWrapper) Kind() reflect.Kind {
	return reflect.TypeOf(p.Payload).Kind()
}

func (p SinglePayloadWrapper) IsKindOf(kind reflect.Kind) bool {
	return p.Kind() == kind
}

func (p SinglePayloadWrapper) Interface() interface{} {
	return p.Payload
}

func (p SinglePayloadWrapper) String() string {
	if p.IsNil() {
		return ""
	}

	if str, ok := p.Payload.(string); ok {
		return str
	}

	return ""
}

func (p SinglePayloadWrapper) SliceString() []string {
	if p.IsNil() {
		return nil
	}

	if arr, ok := p.Payload.([]string); ok {
		return arr
	}

	return nil
}

func (p SinglePayloadWrapper) Boolean() (bool, error) {
	v := p.Payload
	// here we could check for "true", "false" and 0 for false and 1 for true
	// but this may cause unexpected behavior from the developer if they expecting an error
	// so we just check if bool, if yes then return that bool, otherwise return false and an error
	if vb, ok := v.(bool); ok {
		return vb, nil
	}

	if vs, ok := v.(string); ok {
		return strconv.ParseBool(vs)
	}

	return false, fmt.Errorf("unable to parse boolean of %#v", v)
}

func (p SinglePayloadWrapper) Int() (int, error) {
	v := p.Payload

	if vint, ok := v.(int); ok {
		return vint, nil
	}

	if vstring, sok := v.(string); sok {
		return strconv.Atoi(vstring)
	}

	return -1, fmt.Errorf("unable to parse number of %#v", v)
}

func (p SinglePayloadWrapper) Int64() (int64, error) {
	v := p.Payload

	if vint64, ok := v.(int64); ok {
		return vint64, nil
	}

	if vint, ok := v.(int); ok {
		return int64(vint), nil
	}

	if vstring, sok := v.(string); sok {
		return strconv.ParseInt(vstring, 10, 64)
	}

	return -1, fmt.Errorf("unable to parse number of %#v", v)
}

func (p SinglePayloadWrapper) Float32() (float32, error) {
	v := p.Payload

	if vfloat32, ok := v.(float32); ok {
		return vfloat32, nil
	}

	if vfloat64, ok := v.(float64); ok {
		return float32(vfloat64), nil
	}

	if vint, ok := v.(int); ok {
		return float32(vint), nil
	}

	if vstring, sok := v.(string); sok {
		vfloat64, err := strconv.ParseFloat(vstring, 32)
		if err != nil {
			return -1, err
		}
		return float32(vfloat64), nil
	}

	return -1, fmt.Errorf("unable to parse number of %#v", v)
}

func (p SinglePayloadWrapper) Float64() (float64, error) {
	v := p.Payload

	if vfloat32, ok := v.(float32); ok {
		return float64(vfloat32), nil
	}

	if vfloat64, ok := v.(float64); ok {
		return vfloat64, nil
	}

	if vint, ok := v.(int); ok {
		return float64(vint), nil
	}

	if vstring, sok := v.(string); sok {
		return strconv.ParseFloat(vstring, 32)
	}

	return -1, fmt.Errorf("unable to parse number of %#v", v)
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
	payloads []Payload // remember to convert that to simple SInglePayloadWrapper in order to reduce the compiler of creating new object structs on each call.
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

func (p PayloadWrapper) Iterate(visitor func(int, SinglePayloadWrapper)) {
	for i := range p.payloads {
		visitor(i, p.Index(i))
	}
}
