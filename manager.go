// Package ysmrr provides a simple interface for creating and managing
// multiple spinners.
package ysmrr

import (
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chelnak/ysmrr/pkg/charmap"
	"github.com/chelnak/ysmrr/pkg/colors"
	"github.com/chelnak/ysmrr/pkg/tput"
)

// SpinnerManager manages spinners
type SpinnerManager interface {
	AddSpinner(msg string) *Spinner
	GetSpinners() []*Spinner
	GetWriter() io.Writer
	GetCharMap() []string
	GetFrameDuration() time.Duration
	Init()
	Stop()
}

type spinnerManager struct {
	spinners      []*Spinner
	chars         []string
	frameDuration time.Duration
	spinnerColor  colors.Color
	completeColor colors.Color
	errorColor    colors.Color
	messageColor  colors.Color
	writer        io.Writer
	done          chan bool
	ticks         *time.Ticker
	frame         int
}

// AddSpinner adds a new spinner to the manager.
func (sm *spinnerManager) AddSpinner(message string) *Spinner {
	opts := SpinnerOptions{
		Message:       message,
		SpinnerColor:  sm.spinnerColor,
		CompleteColor: sm.completeColor,
		ErrorColor:    sm.errorColor,
		MessageColor:  sm.messageColor,
	}

	spinner := NewSpinner(opts)
	sm.spinners = append(sm.spinners, spinner)
	return spinner
}

// GetSpinners returns the spinners managed by the manager.
func (sm *spinnerManager) GetSpinners() []*Spinner {
	return sm.spinners
}

// Init initializes the spinnerManager and starts the renderer.
func (sm *spinnerManager) Init() {
	// Handle SIGINT and SIGTERM so we can ensure that the
	// terminal is properly reset.
	// Unsure if this is the right place for this especially given
	// that it calls os.Exit.
	signals := make(chan os.Signal, 2)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		sm.Stop()
		os.Exit(0)
	}()

	sm.ticks = time.NewTicker(sm.frameDuration)
	go sm.render()
}

// Stop signals that all spinners should complete.
func (sm *spinnerManager) Stop() {
	sm.done <- true
	sm.ticks.Stop()
	defer tput.Cnorm(sm.writer)

	// Persist the final frame for each spinner.
	for _, s := range sm.spinners {
		s.Print(sm.writer, sm.chars[sm.frame])
	}
}

// GetWriter returns the configured io.Writer.
func (sm *spinnerManager) GetWriter() io.Writer {
	return sm.writer
}

// GetCharMap returns the configured character map.
func (sm *spinnerManager) GetCharMap() []string {
	return sm.chars
}

// GetFrameDuration returns the configured frame duration.
func (sm *spinnerManager) GetFrameDuration() time.Duration {
	return sm.frameDuration
}

func (sm *spinnerManager) setNextFrame() {
	sm.frame += 1
	if sm.frame >= len(sm.chars) {
		sm.frame = 0
	}
}

func (sm *spinnerManager) renderFrame() {
	for _, s := range sm.spinners {
		s.Print(sm.writer, sm.chars[sm.frame])
	}
	sm.setNextFrame()
}

func (sm *spinnerManager) render() {
	tput.Sc(sm.writer)
	tput.Civis(sm.writer)

	for {
		select {
		case <-sm.done:
			return
		case <-sm.ticks.C:
			sm.renderFrame()
		}

		tput.Rc(sm.writer)
	}
}

// Option represents a spinner manager option.
type Option func(*spinnerManager)

// NewSpinnerManager creates a new spinner manager.
func NewSpinnerManager(options ...Option) SpinnerManager {
	sm := &spinnerManager{
		chars:         charmap.Dots,
		frameDuration: 250 * time.Millisecond,
		spinnerColor:  colors.FgHiGreen,
		errorColor:    colors.FgHiRed,
		completeColor: colors.FgHiGreen,
		messageColor:  colors.NoColor,
		writer:        os.Stdout,
		done:          make(chan bool),
	}

	for _, option := range options {
		option(sm)
	}

	return sm
}

// WithCharMap sets the characters used for the spinners.
// Available charmaps can be found in the package github.com/chelnak/ysmrr/pkg/charmap.
// The default charmap is the Dots.
func WithCharMap(chars []string) Option {
	return func(sm *spinnerManager) {
		sm.chars = chars
	}
}

// WithFrameDuration sets the duration of each frame.
// The default duration is 250 milliseconds.
func WithFrameDuration(d time.Duration) Option {
	return func(sm *spinnerManager) {
		sm.frameDuration = d
	}
}

// WithSpinnerColor sets the color of the spinners.
// Available colors can be found in the package github.com/chelnak/ysmrr/pkg/colors.
// The default color is FgHiGreen.
func WithSpinnerColor(c colors.Color) Option {
	return func(sm *spinnerManager) {
		sm.spinnerColor = c
	}
}

// WithErrorColor sets the color of the error icon.
// Available colors can be found in the package github.com/chelnak/ysmrr/pkg/colors.
// The default color is FgHiRed.
func WithErrorColor(c colors.Color) Option {
	return func(sm *spinnerManager) {
		sm.errorColor = c
	}
}

// WithCompleteColor sets the color of the complete icon.
// Available colors can be found in the package github.com/chelnak/ysmrr/pkg/colors.
// The default color is FgHiGreen.
func WithCompleteColor(c colors.Color) Option {
	return func(sm *spinnerManager) {
		sm.completeColor = c
	}
}

// WithMessageColor sets the color of the message.
// Available colors can be found in the package github.com/chelnak/ysmrr/pkg/colors.
// The default color is NoColor.
func WithMessageColor(c colors.Color) Option {
	return func(sm *spinnerManager) {
		sm.messageColor = c
	}
}

// WithWriter sets the writer used for the spinners.
// The writer can be anything that implements the io.Writer interface.
// The default writer is os.Stdout.
func WithWriter(w io.Writer) Option {
	return func(sm *spinnerManager) {
		sm.writer = w
	}
}
