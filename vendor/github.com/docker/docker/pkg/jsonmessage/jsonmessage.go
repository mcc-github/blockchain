package jsonmessage 

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	gotty "github.com/Nvveen/Gotty"
	"github.com/docker/docker/pkg/term"
	units "github.com/docker/go-units"
)



const RFC3339NanoFixed = "2006-01-02T15:04:05.000000000Z07:00"



type JSONError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *JSONError) Error() string {
	return e.Message
}





type JSONProgress struct {
	terminalFd uintptr
	Current    int64 `json:"current,omitempty"`
	Total      int64 `json:"total,omitempty"`
	Start      int64 `json:"start,omitempty"`
	
	HideCounts bool   `json:"hidecounts,omitempty"`
	Units      string `json:"units,omitempty"`
	nowFunc    func() time.Time
	winSize    int
}

func (p *JSONProgress) String() string {
	var (
		width       = p.width()
		pbBox       string
		numbersBox  string
		timeLeftBox string
	)
	if p.Current <= 0 && p.Total <= 0 {
		return ""
	}
	if p.Total <= 0 {
		switch p.Units {
		case "":
			current := units.HumanSize(float64(p.Current))
			return fmt.Sprintf("%8v", current)
		default:
			return fmt.Sprintf("%d %s", p.Current, p.Units)
		}
	}

	percentage := int(float64(p.Current)/float64(p.Total)*100) / 2
	if percentage > 50 {
		percentage = 50
	}
	if width > 110 {
		
		numSpaces := 0
		if 50-percentage > 0 {
			numSpaces = 50 - percentage
		}
		pbBox = fmt.Sprintf("[%s>%s] ", strings.Repeat("=", percentage), strings.Repeat(" ", numSpaces))
	}

	switch {
	case p.HideCounts:
	case p.Units == "": 
		current := units.HumanSize(float64(p.Current))
		total := units.HumanSize(float64(p.Total))

		numbersBox = fmt.Sprintf("%8v/%v", current, total)

		if p.Current > p.Total {
			
			numbersBox = fmt.Sprintf("%8v", current)
		}
	default:
		numbersBox = fmt.Sprintf("%d/%d %s", p.Current, p.Total, p.Units)

		if p.Current > p.Total {
			
			numbersBox = fmt.Sprintf("%d %s", p.Current, p.Units)
		}
	}

	if p.Current > 0 && p.Start > 0 && percentage < 50 {
		fromStart := p.now().Sub(time.Unix(p.Start, 0))
		perEntry := fromStart / time.Duration(p.Current)
		left := time.Duration(p.Total-p.Current) * perEntry
		left = (left / time.Second) * time.Second

		if width > 50 {
			timeLeftBox = " " + left.String()
		}
	}
	return pbBox + numbersBox + timeLeftBox
}


func (p *JSONProgress) now() time.Time {
	if p.nowFunc == nil {
		p.nowFunc = func() time.Time {
			return time.Now().UTC()
		}
	}
	return p.nowFunc()
}


func (p *JSONProgress) width() int {
	if p.winSize != 0 {
		return p.winSize
	}
	ws, err := term.GetWinsize(p.terminalFd)
	if err == nil {
		return int(ws.Width)
	}
	return 200
}




type JSONMessage struct {
	Stream          string        `json:"stream,omitempty"`
	Status          string        `json:"status,omitempty"`
	Progress        *JSONProgress `json:"progressDetail,omitempty"`
	ProgressMessage string        `json:"progress,omitempty"` 
	ID              string        `json:"id,omitempty"`
	From            string        `json:"from,omitempty"`
	Time            int64         `json:"time,omitempty"`
	TimeNano        int64         `json:"timeNano,omitempty"`
	Error           *JSONError    `json:"errorDetail,omitempty"`
	ErrorMessage    string        `json:"error,omitempty"` 
	
	Aux *json.RawMessage `json:"aux,omitempty"`
}


type termInfo interface {
	Parse(attr string, params ...interface{}) (string, error)
}

type noTermInfo struct{} 

func (ti *noTermInfo) Parse(attr string, params ...interface{}) (string, error) {
	return "", fmt.Errorf("noTermInfo")
}

func clearLine(out io.Writer, ti termInfo) {
	

	
	if attr, err := ti.Parse("el1"); err == nil {
		fmt.Fprintf(out, "%s", attr)
	} else {
		fmt.Fprintf(out, "\x1b[1K")
	}
	
	if attr, err := ti.Parse("el"); err == nil {
		fmt.Fprintf(out, "%s", attr)
	} else {
		fmt.Fprintf(out, "\x1b[K")
	}
}

func cursorUp(out io.Writer, ti termInfo, l int) {
	if l == 0 { 
		return
	}
	if attr, err := ti.Parse("cuu", l); err == nil {
		fmt.Fprintf(out, "%s", attr)
	} else {
		fmt.Fprintf(out, "\x1b[%dA", l)
	}
}

func cursorDown(out io.Writer, ti termInfo, l int) {
	if l == 0 { 
		return
	}
	if attr, err := ti.Parse("cud", l); err == nil {
		fmt.Fprintf(out, "%s", attr)
	} else {
		fmt.Fprintf(out, "\x1b[%dB", l)
	}
}




func (jm *JSONMessage) Display(out io.Writer, termInfo termInfo) error {
	if jm.Error != nil {
		if jm.Error.Code == 401 {
			return fmt.Errorf("authentication is required")
		}
		return jm.Error
	}
	var endl string
	if termInfo != nil && jm.Stream == "" && jm.Progress != nil {
		clearLine(out, termInfo)
		endl = "\r"
		fmt.Fprintf(out, endl)
	} else if jm.Progress != nil && jm.Progress.String() != "" { 
		return nil
	}
	if jm.TimeNano != 0 {
		fmt.Fprintf(out, "%s ", time.Unix(0, jm.TimeNano).Format(RFC3339NanoFixed))
	} else if jm.Time != 0 {
		fmt.Fprintf(out, "%s ", time.Unix(jm.Time, 0).Format(RFC3339NanoFixed))
	}
	if jm.ID != "" {
		fmt.Fprintf(out, "%s: ", jm.ID)
	}
	if jm.From != "" {
		fmt.Fprintf(out, "(from %s) ", jm.From)
	}
	if jm.Progress != nil && termInfo != nil {
		fmt.Fprintf(out, "%s %s%s", jm.Status, jm.Progress.String(), endl)
	} else if jm.ProgressMessage != "" { 
		fmt.Fprintf(out, "%s %s%s", jm.Status, jm.ProgressMessage, endl)
	} else if jm.Stream != "" {
		fmt.Fprintf(out, "%s%s", jm.Stream, endl)
	} else {
		fmt.Fprintf(out, "%s%s\n", jm.Status, endl)
	}
	return nil
}




func DisplayJSONMessagesStream(in io.Reader, out io.Writer, terminalFd uintptr, isTerminal bool, auxCallback func(*json.RawMessage)) error {
	var (
		dec = json.NewDecoder(in)
		ids = make(map[string]int)
	)

	var termInfo termInfo

	if isTerminal {
		term := os.Getenv("TERM")
		if term == "" {
			term = "vt102"
		}

		var err error
		if termInfo, err = gotty.OpenTermInfo(term); err != nil {
			termInfo = &noTermInfo{}
		}
	}

	for {
		diff := 0
		var jm JSONMessage
		if err := dec.Decode(&jm); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if jm.Aux != nil {
			if auxCallback != nil {
				auxCallback(jm.Aux)
			}
			continue
		}

		if jm.Progress != nil {
			jm.Progress.terminalFd = terminalFd
		}
		if jm.ID != "" && (jm.Progress != nil || jm.ProgressMessage != "") {
			line, ok := ids[jm.ID]
			if !ok {
				
				
				
				
				
				
				line = len(ids)
				ids[jm.ID] = line
				if termInfo != nil {
					fmt.Fprintf(out, "\n")
				}
			}
			diff = len(ids) - line
			if termInfo != nil {
				cursorUp(out, termInfo, diff)
			}
		} else {
			
			
			
			
			
			ids = make(map[string]int)
		}
		err := jm.Display(out, termInfo)
		if jm.ID != "" && termInfo != nil {
			cursorDown(out, termInfo, diff)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

type stream interface {
	io.Writer
	FD() uintptr
	IsTerminal() bool
}


func DisplayJSONMessagesToStream(in io.Reader, stream stream, auxCallback func(*json.RawMessage)) error {
	return DisplayJSONMessagesStream(in, stream, stream.FD(), stream.IsTerminal(), auxCallback)
}
