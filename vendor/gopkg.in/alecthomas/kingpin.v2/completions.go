package kingpin



type HintAction func() []string
type completionsMixin struct {
	hintActions        []HintAction
	builtinHintActions []HintAction
}

func (a *completionsMixin) addHintAction(action HintAction) {
	a.hintActions = append(a.hintActions, action)
}


func (a *completionsMixin) addHintActionBuiltin(action HintAction) {
	a.builtinHintActions = append(a.builtinHintActions, action)
}

func (a *completionsMixin) resolveCompletions() []string {
	var hints []string

	options := a.builtinHintActions
	if len(a.hintActions) > 0 {
		
		options = a.hintActions
	}

	for _, hintAction := range options {
		hints = append(hints, hintAction()...)
	}
	return hints
}
