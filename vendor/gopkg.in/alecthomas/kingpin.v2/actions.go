package kingpin




type Action func(*ParseContext) error

type actionMixin struct {
	actions    []Action
	preActions []Action
}

type actionApplier interface {
	applyActions(*ParseContext) error
	applyPreActions(*ParseContext) error
}

func (a *actionMixin) addAction(action Action) {
	a.actions = append(a.actions, action)
}

func (a *actionMixin) addPreAction(action Action) {
	a.preActions = append(a.preActions, action)
}

func (a *actionMixin) applyActions(context *ParseContext) error {
	for _, action := range a.actions {
		if err := action(context); err != nil {
			return err
		}
	}
	return nil
}

func (a *actionMixin) applyPreActions(context *ParseContext) error {
	for _, preAction := range a.preActions {
		if err := preAction(context); err != nil {
			return err
		}
	}
	return nil
}
