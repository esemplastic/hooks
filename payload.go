package hooks

import (
	"fmt"
	"net"
	"net/http"
	"reflect"
	"strconv"
)

///TODO: commenting and docs when ready.

func ReadPayload(data interface{}) Payload {
	return payload{data}
}

// Payload is an interface in order to be able to change its methods and use a custom Payload when needed.
type Payload interface {
	IsNil() bool
	Kind() reflect.Kind
	IsKindOf(kind reflect.Kind) bool
	Interface() interface{}
	String() string
	SliceString() []string
	Len() int
	Boolean() (bool, error)
	Int() (int, error)
	Int64() (int64, error)
	Float32() (float32, error)
	Float64() (float64, error)
	Err() error
	Listener() net.Listener
	Handler() http.Handler
}

type payload struct {
	payload interface{}
}

func (p payload) IsNil() bool {
	return p.payload == nil
}

func (p payload) Kind() reflect.Kind {
	return reflect.TypeOf(p.payload).Kind()
}

func (p payload) IsKindOf(kind reflect.Kind) bool {
	return p.Kind() == kind
}

func (p payload) Interface() interface{} {
	return p.payload
}

func (p payload) String() string {
	if p.IsNil() {
		return ""
	}

	if str, ok := p.payload.(string); ok {
		return str
	}

	return ""
}

func (p payload) SliceString() []string {
	if p.IsNil() {
		return nil
	}

	if arr, ok := p.payload.([]string); ok {
		return arr
	}

	return nil
}

func (p payload) Len() int {
	if typ := reflect.TypeOf(p.payload); typ.Kind() == reflect.Slice {
		return typ.Len()
	}

	return 0
}

func (p payload) Boolean() (bool, error) {
	v := p.payload
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

func (p payload) Int() (int, error) {
	v := p.payload

	if vint, ok := v.(int); ok {
		return vint, nil
	}

	if vstring, sok := v.(string); sok {
		return strconv.Atoi(vstring)
	}

	return -1, fmt.Errorf("unable to parse number of %#v", v)
}

func (p payload) Int64() (int64, error) {
	v := p.payload

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

func (p payload) Float32() (float32, error) {
	v := p.payload

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

func (p payload) Float64() (float64, error) {
	v := p.payload

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

func (p payload) Err() error {
	if err, ok := p.payload.(error); ok {
		return err
	}

	return nil
}

func (p payload) Listener() net.Listener {
	if ln, ok := p.payload.(net.Listener); ok {
		return ln
	}

	return nil
}

func (p payload) Handler() http.Handler {
	if h, ok := p.payload.(http.Handler); ok {
		return h
	}

	return nil
}

type Payloads interface {
	Index(idx int) Payload
	First() Payload
	Second() Payload
	Last() Payload
	Range(start int, end int) Payloads
	Iterate(visitor func(int, Payload))
	Len() int
}

type payloads []Payload

func (p payloads) Index(idx int) Payload {
	if idx+1 > len(p) {
		return payload{}
	}

	return p[idx]
}

func (p payloads) First() Payload {
	return p.Index(0)
}

func (p payloads) Second() Payload {
	return p.Index(1)
}

func (p payloads) Last() Payload {
	return p.Index(len(p) - 1)
}

// including start, excluding end
func (p payloads) Range(start int, end int) Payloads {
	if end >= len(p) {
		return payloads{}

	}
	return p[start:end]
}

func (p payloads) Iterate(visitor func(int, Payload)) {
	for i := range p {
		visitor(i, p.Index(i))
	}
}

func (p payloads) Len() int {
	return len(p)
}
