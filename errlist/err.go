package errlist

import (
	"fmt"
)

type ErrNode struct {
	Data map[string]string
	err  error
	next *ErrNode
}

func New(err error) (self *ErrNode) {
	if errAsErr, ok := err.(*ErrNode); ok {
		return errAsErr
	}

	// to prevent segfault on Unwrap().Error() of childless node with this error inside
	if err == nil {
		err = ErrEmpty
	}

	return &ErrNode{
		Data: make(map[string]string),
		err:  err,
	}
}

// Checks if `e` is standalone node or has children.
func (e *ErrNode) HasChildren() bool {
	return e.next == nil
}

// Sets data inside underlying map at `k`.
func (e *ErrNode) Set(k, v string) (self *ErrNode) {
	e.Data[k] = v
	return e
}

// Gets data from underlying map at `k`.
func (e *ErrNode) Get(k string) (v string, ok bool) {
	v, ok = e.Data[k]
	return v, ok
}

// `e` wraps `child`. If `child` is not of type `ErrNode`, `New()` is called.
func (e *ErrNode) Wrap(child error) (self *ErrNode) {
	tail := e
	for tail.next != nil {
		tail = tail.next
	}

	if childAsErr, ok := child.(*ErrNode); ok {
		tail.next = childAsErr
		return e
	}
	tail.next = New(child)
	return e
}

// If `e.next` is not `nil` returns `next` while acting like a list pop back.
// Otherwise returs underlying `error`.
func (e *ErrNode) Unwrap() error {
	tail := e
	if tail.next == nil {
		return tail.err
	}

	var prev *ErrNode
	for tail.next != nil {
		prev = tail
		tail = tail.next
	}

	res := *tail
	prev.next = nil

	return &res
}

// Same as `Unwrap()` but returns self when called on childless node.
func (e *ErrNode) UnwrapAsNode() *ErrNode {
	tail := e
	if tail.next == nil {
		return tail
	}

	var prev *ErrNode
	for tail.next != nil {
		prev = tail
		tail = tail.next
	}

	res := *tail
	prev.next = nil

	return &res
}

// Rerurns `e`'s represented as json string.
func (e *ErrNode) json() string {
	var res string

	if e.err != ErrEmpty {
		res = fmt.Sprintf("\"error\": \"%v\"", e.err)
	}

	if len(e.Data) > 0 {
		var data string
		for k, v := range e.Data {
			data += fmt.Sprintf("\"%s\": \"%s\", ", k, v)
		}
		data = fmt.Sprintf("{%s}", data[:len(data)-2])
		if res != "" {
			res = fmt.Sprintf("%s, \"data\": %s", res, data)
		} else {
			res = fmt.Sprintf("\"data\": %s", data)
		}
	}
	return fmt.Sprintf("{%s}", res)
}

func (e ErrNode) Error() string {
	res := e.json() + "\n"
	err := e
	depth := 0
	for err.next != nil {
		for i := 0; i < depth; i++ {
			res += "    "
		}
		depth++

		res += "  L " + err.next.json() + "\n"
		err = *err.next
	}

	return res
}

func (e *ErrNode) JSON() string {
	/* "errors": [
		{...},
		{...},
		...
	] */

	res := "\"errors\": [\n" + "    " + e.json()
	err := e
	for err.next != nil {
		res += ",\n" + "    " + err.next.json()
		err = err.next
	}
	res += "\n]"

	return res

}
