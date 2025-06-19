package debugger

import "fmt"

type Debugger struct {
	enabled bool
}

func New(enabled bool) *Debugger {
	return &Debugger{enabled}
}

func (d *Debugger) Log(a ...any) {
	if d.enabled {
		fmt.Println(a...)
	}
}

func (d *Debugger) Logf(format string, a ...any) {
	if d.enabled {
		fmt.Println(fmt.Sprintf(format, a...))
	}
}
