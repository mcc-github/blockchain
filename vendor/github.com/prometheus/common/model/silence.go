












package model

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"
)


type Matcher struct {
	Name    LabelName `json:"name"`
	Value   string    `json:"value"`
	IsRegex bool      `json:"isRegex"`
}

func (m *Matcher) UnmarshalJSON(b []byte) error {
	type plain Matcher
	if err := json.Unmarshal(b, (*plain)(m)); err != nil {
		return err
	}

	if len(m.Name) == 0 {
		return fmt.Errorf("label name in matcher must not be empty")
	}
	if m.IsRegex {
		if _, err := regexp.Compile(m.Value); err != nil {
			return err
		}
	}
	return nil
}


func (m *Matcher) Validate() error {
	if !m.Name.IsValid() {
		return fmt.Errorf("invalid name %q", m.Name)
	}
	if m.IsRegex {
		if _, err := regexp.Compile(m.Value); err != nil {
			return fmt.Errorf("invalid regular expression %q", m.Value)
		}
	} else if !LabelValue(m.Value).IsValid() || len(m.Value) == 0 {
		return fmt.Errorf("invalid value %q", m.Value)
	}
	return nil
}



type Silence struct {
	ID uint64 `json:"id,omitempty"`

	Matchers []*Matcher `json:"matchers"`

	StartsAt time.Time `json:"startsAt"`
	EndsAt   time.Time `json:"endsAt"`

	CreatedAt time.Time `json:"createdAt,omitempty"`
	CreatedBy string    `json:"createdBy"`
	Comment   string    `json:"comment,omitempty"`
}


func (s *Silence) Validate() error {
	if len(s.Matchers) == 0 {
		return fmt.Errorf("at least one matcher required")
	}
	for _, m := range s.Matchers {
		if err := m.Validate(); err != nil {
			return fmt.Errorf("invalid matcher: %s", err)
		}
	}
	if s.StartsAt.IsZero() {
		return fmt.Errorf("start time missing")
	}
	if s.EndsAt.IsZero() {
		return fmt.Errorf("end time missing")
	}
	if s.EndsAt.Before(s.StartsAt) {
		return fmt.Errorf("start time must be before end time")
	}
	if s.CreatedBy == "" {
		return fmt.Errorf("creator information missing")
	}
	if s.Comment == "" {
		return fmt.Errorf("comment missing")
	}
	if s.CreatedAt.IsZero() {
		return fmt.Errorf("creation timestamp missing")
	}
	return nil
}
