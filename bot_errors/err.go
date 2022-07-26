package bot_errors

type Err struct {
	session string
	event   string
	err     error
	next    *Err
}

func New(session, event string, err error) *Err {
	return &Err{
		session: session,
		event:   event,
		err:     err,
	}
}

func (e *Err) Nest(child error) *Err {
	tail := e
	for tail.next != nil {
		tail = tail.next
	}

	if childAsErr, ok := child.(*Err); ok {
		tail.next = childAsErr
		return e
	}

	tail.next = New("", "", child)
	return e
}

func (e *Err) strln() string {
	return e.session + " " + e.event + " " + e.err.Error() + "\n"
}

func (e *Err) Unwrap() error { // test
	tail := e
	if tail.next == nil {
		return tail
	}

	var prev *Err
	for tail.next != nil {
		prev = tail
		tail = tail.next
	}

	res := *tail
	prev.next = nil

	return &res
}

func (e Err) Error() string {
	res := e.strln()
	err := e
	depth := 0
	for err.next != nil {
		for i := 0; i < depth; i++ {
			res += "    "
		}
		depth++

		res += "  L " + err.next.strln()
		err = *err.next
	}

	return res
}
