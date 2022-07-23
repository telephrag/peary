package bot_errors

import (
	"fmt"
)

type Err struct {
	Session string // doubles as session id from taking role to losing it by one way or another
	Event   string
	Next    *Nested
}

func (e *Err) Nest(n *Nested) {
	tail := e.Next
	for tail != nil {
		tail = tail.Next
	}
	tail = n
}

func (e *Err) String() string {
	res := e.Session + " " + e.Event + "\n"
	nested := e.Next
	depth := 1
	for nested != nil {
		for i := 0; i < depth; i++ {
			res += "  "
		}

		res += "L "

		res = fmt.Sprintf("%s%v\n", res, nested.Err)
	}
	return res
}

type Nested struct {
	Event string
	Err   error
	Next  *Nested
}
