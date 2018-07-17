package ansiterm

type oscStringState struct {
	baseState
}

func (oscState oscStringState) Handle(b byte) (s state, e error) {
	oscState.parser.logf("OscString::Handle %#x", b)
	nextState, err := oscState.baseState.Handle(b)
	if nextState != nil || err != nil {
		return nextState, err
	}

	switch {
	case isOscStringTerminator(b):
		return oscState.parser.ground, nil
	}

	return oscState, nil
}



func isOscStringTerminator(b byte) bool {

	if b == ANSI_BEL || b == 0x5C {
		return true
	}

	return false
}
