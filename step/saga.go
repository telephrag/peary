package step

type Saga struct {
	iterPos int
	steps   []Step
}

func NewSaga(steps []Step) *Saga {
	return &Saga{
		iterPos: -1,
		steps:   steps,
	}
}

func (s *Saga) GetStep() Step {
	if s.iterPos < 0 || s.iterPos > len(s.steps)-1 {
		return nil
	}
	return s.steps[s.iterPos]
}

func (s *Saga) Next() Step {
	if s.iterPos+1 < len(s.steps) {
		s.iterPos++
		return s.steps[s.iterPos]
	}
	return nil
}

func (s *Saga) Prev() Step {
	if s.iterPos-1 >= 0 {
		s.iterPos--
		return s.steps[s.iterPos]
	}
	return nil
}

func (s *Saga) ResetIter() {
	s.iterPos = -1
}
