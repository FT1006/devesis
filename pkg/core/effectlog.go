package core

import (
	"fmt"
	"time"
)

// EffectLog captures all game state mutations for preview/resolve display
type EffectLog struct {
	Lines []string
}

// Add appends a formatted log line to the effect log
func (l *EffectLog) Add(msg string, args ...any) {
	l.Lines = append(l.Lines, fmt.Sprintf(msg, args...))
}

// Clear empties the log for reuse
func (l *EffectLog) Clear() {
	l.Lines = l.Lines[:0]
}

// IsEmpty returns true if no effects have been logged
func (l *EffectLog) IsEmpty() bool {
	return len(l.Lines) == 0
}

// StreamLines prints each log line with a small delay for step-by-step feel
func (l *EffectLog) StreamLines(delay time.Duration) {
	for _, line := range l.Lines {
		fmt.Println(line)
		if delay > 0 {
			time.Sleep(delay)
		}
	}
}

// PrintBulk prints all lines at once without delay
func (l *EffectLog) PrintBulk() {
	for _, line := range l.Lines {
		fmt.Println(line)
	}
}

// NewEffectLog creates a new empty effect log
func NewEffectLog() *EffectLog {
	return &EffectLog{
		Lines: make([]string, 0),
	}
}