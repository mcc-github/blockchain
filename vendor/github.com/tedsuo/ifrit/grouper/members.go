package grouper

import (
	"fmt"

	"github.com/tedsuo/ifrit"
)


type Member struct {
	Name string
	ifrit.Runner
}


type Members []Member


func (m Members) Validate() error {
	foundNames := map[string]struct{}{}
	foundToken := struct{}{}
	duplicateNames := []string{}

	for _, member := range m {
		_, present := foundNames[member.Name]
		if present {
			duplicateNames = append(duplicateNames, member.Name)
			continue
		}
		foundNames[member.Name] = foundToken
	}

	if len(duplicateNames) > 0 {
		return ErrDuplicateNames{duplicateNames}
	}
	return nil
}


type ErrDuplicateNames struct {
	DuplicateNames []string
}

func (e ErrDuplicateNames) Error() string {
	var msg string

	switch len(e.DuplicateNames) {
	case 0:
		msg = fmt.Sprintln("ErrDuplicateNames initialized without any duplicate names.")
	case 1:
		msg = fmt.Sprintln("Duplicate member name:", e.DuplicateNames[0])
	default:
		msg = fmt.Sprintln("Duplicate member names:")
		for _, name := range e.DuplicateNames {
			msg = fmt.Sprintln(name)
		}
	}

	return msg
}
