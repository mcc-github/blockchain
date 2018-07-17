












package model

import (
	"fmt"
	"time"
)

type AlertStatus string

const (
	AlertFiring   AlertStatus = "firing"
	AlertResolved AlertStatus = "resolved"
)


type Alert struct {
	
	
	Labels LabelSet `json:"labels"`

	
	Annotations LabelSet `json:"annotations"`

	
	StartsAt     time.Time `json:"startsAt,omitempty"`
	EndsAt       time.Time `json:"endsAt,omitempty"`
	GeneratorURL string    `json:"generatorURL"`
}


func (a *Alert) Name() string {
	return string(a.Labels[AlertNameLabel])
}



func (a *Alert) Fingerprint() Fingerprint {
	return a.Labels.Fingerprint()
}

func (a *Alert) String() string {
	s := fmt.Sprintf("%s[%s]", a.Name(), a.Fingerprint().String()[:7])
	if a.Resolved() {
		return s + "[resolved]"
	}
	return s + "[active]"
}


func (a *Alert) Resolved() bool {
	return a.ResolvedAt(time.Now())
}



func (a *Alert) ResolvedAt(ts time.Time) bool {
	if a.EndsAt.IsZero() {
		return false
	}
	return !a.EndsAt.After(ts)
}


func (a *Alert) Status() AlertStatus {
	if a.Resolved() {
		return AlertResolved
	}
	return AlertFiring
}


func (a *Alert) Validate() error {
	if a.StartsAt.IsZero() {
		return fmt.Errorf("start time missing")
	}
	if !a.EndsAt.IsZero() && a.EndsAt.Before(a.StartsAt) {
		return fmt.Errorf("start time must be before end time")
	}
	if err := a.Labels.Validate(); err != nil {
		return fmt.Errorf("invalid label set: %s", err)
	}
	if len(a.Labels) == 0 {
		return fmt.Errorf("at least one label pair required")
	}
	if err := a.Annotations.Validate(); err != nil {
		return fmt.Errorf("invalid annotations: %s", err)
	}
	return nil
}


type Alerts []*Alert

func (as Alerts) Len() int      { return len(as) }
func (as Alerts) Swap(i, j int) { as[i], as[j] = as[j], as[i] }

func (as Alerts) Less(i, j int) bool {
	if as[i].StartsAt.Before(as[j].StartsAt) {
		return true
	}
	if as[i].EndsAt.Before(as[j].EndsAt) {
		return true
	}
	return as[i].Fingerprint() < as[j].Fingerprint()
}


func (as Alerts) HasFiring() bool {
	for _, a := range as {
		if !a.Resolved() {
			return true
		}
	}
	return false
}


func (as Alerts) Status() AlertStatus {
	if as.HasFiring() {
		return AlertFiring
	}
	return AlertResolved
}
