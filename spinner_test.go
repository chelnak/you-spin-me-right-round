package ysmrr_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/chelnak/ysmrr"
	"github.com/chelnak/ysmrr/pkg/animations"
	"github.com/chelnak/ysmrr/pkg/colors"
	"github.com/stretchr/testify/assert"
)

var initialMessage = "test"
var initialOpts = ysmrr.SpinnerOptions{
	Message:           initialMessage,
	SpinnerColor:      colors.NoColor,
	CompleteColor:     colors.NoColor,
	ErrorColor:        colors.NoColor,
	MessageColor:      colors.NoColor,
	CompleteCharacter: "✓",
	ErrorCharacter:    "✗",
	HasUpdate:         make(chan bool),
}

func TestNewSpinner(t *testing.T) {
	opts := initialOpts
	spinner := ysmrr.NewSpinner(opts)
	assert.NotNil(t, spinner)
}

func TestSpinnerGetMessage(t *testing.T) {
	opts := initialOpts
	spinner := ysmrr.NewSpinner(opts)
	assert.Equal(t, initialMessage, spinner.GetMessage())
}

func TestSpinnerGetPrefix(t *testing.T) {
	opts := initialOpts
	spinner := ysmrr.NewSpinner(opts)
	assert.Equal(t, "", spinner.GetPrefix())
}

func TestSpinnerIsError(t *testing.T) {
	opts := initialOpts
	spinner := ysmrr.NewSpinner(opts)
	assert.Equal(t, false, spinner.IsError())
}

func TestSpinnerIsComplete(t *testing.T) {
	opts := initialOpts
	spinner := ysmrr.NewSpinner(opts)
	assert.Equal(t, false, spinner.IsComplete())
}

func TestSpinnerUpdateMessage(t *testing.T) {
	updatedMessage := "updated message"
	opts := initialOpts
	spinner := ysmrr.NewSpinner(opts)
	spinner.UpdateMessage(updatedMessage)
	assert.Equal(t, updatedMessage, spinner.GetMessage())
}

func TestSpinnerUpdateMessagef(t *testing.T) {
	expectedMessage := "updated message test"
	opts := initialOpts
	spinner := ysmrr.NewSpinner(opts)
	spinner.UpdateMessagef("updated message %s", "test")
	assert.Equal(t, expectedMessage, spinner.GetMessage())
}

func TestSpinnerUpdatePrefix(t *testing.T) {
	expectedPrefix := "prefix"
	opts := initialOpts
	spinner := ysmrr.NewSpinner(opts)
	spinner.UpdatePrefix(expectedPrefix)
	assert.Equal(t, expectedPrefix, spinner.GetPrefix())
}

func TestSpinnerUpdatePrefixf(t *testing.T) {
	expectedPrefix := "prefix test"
	opts := initialOpts
	spinner := ysmrr.NewSpinner(opts)
	spinner.UpdatePrefixf("prefix %s", "test")
	assert.Equal(t, expectedPrefix, spinner.GetPrefix())
}

func TestSpinnerCompleteWithMessage(t *testing.T) {
	opts := initialOpts
	spinner := ysmrr.NewSpinner(opts)
	spinner.CompleteWithMessage("complete")
	assert.Equal(t, true, spinner.IsComplete())
	assert.Equal(t, "complete", spinner.GetMessage())
}

func TestSpinnerCompleteWithMessagef(t *testing.T) {
	opts := initialOpts
	spinner := ysmrr.NewSpinner(opts)
	spinner.CompleteWithMessagef("complete %s", "test")
	assert.Equal(t, true, spinner.IsComplete())
	assert.Equal(t, "complete test", spinner.GetMessage())
}

func TestSpinnerComplete(t *testing.T) {
	opts := initialOpts
	spinner := ysmrr.NewSpinner(opts)
	spinner.Complete()
	assert.Equal(t, true, spinner.IsComplete())
}

func TestCompleteCharacter(t *testing.T) {
	opts := initialOpts
	expectedCharacter := "*"
	spinner := ysmrr.NewSpinner(opts)
	spinner.CompleteCharacter(expectedCharacter)
	assert.Equal(t, expectedCharacter, spinner.GetCompleteCharacter())
}

func TestErrorCharacter(t *testing.T) {
	opts := initialOpts
	expectedCharacter := "*"
	spinner := ysmrr.NewSpinner(opts)
	spinner.ErrorCharacter(expectedCharacter)
	assert.Equal(t, expectedCharacter, spinner.GetErrorCharacter())
}

func TestSpinnerErrorWithMessage(t *testing.T) {
	opts := initialOpts
	spinner := ysmrr.NewSpinner(opts)
	spinner.ErrorWithMessage("error")
	assert.Equal(t, true, spinner.IsError())
	assert.Equal(t, "error", spinner.GetMessage())
}

func TestSpinnerErrorWithMessagef(t *testing.T) {
	opts := initialOpts
	spinner := ysmrr.NewSpinner(opts)
	spinner.ErrorWithMessagef("error %s", "test")
	assert.Equal(t, true, spinner.IsError())
	assert.Equal(t, "error test", spinner.GetMessage())
}

func TestSpinnerError(t *testing.T) {
	opts := initialOpts
	spinner := ysmrr.NewSpinner(opts)
	spinner.Error()
	assert.Equal(t, true, spinner.IsError())
}

func TestPrint(t *testing.T) {
	opts := initialOpts
	spinner := ysmrr.NewSpinner(opts)

	var buf bytes.Buffer
	_, dots := animations.GetAnimation(animations.Dots)
	spinner.Print(&buf, dots[0])

	want := fmt.Sprintf("%s %s\r\n", dots[0], initialMessage)
	assert.Equal(t, want, buf.String())
}

func TestPrintWithPrefix(t *testing.T) {
	opts := initialOpts
	spinner := ysmrr.NewSpinner(opts)

	prefix := "prefix"
	spinner.UpdatePrefix(prefix)

	var buf bytes.Buffer
	_, dots := animations.GetAnimation(animations.Dots)
	spinner.Print(&buf, dots[0])

	want := fmt.Sprintf("%s%s %s\r\n", prefix, dots[0], initialMessage)
	assert.Equal(t, want, buf.String())
}

func TestPrintWithComplete(t *testing.T) {
	opts := initialOpts
	spinner := ysmrr.NewSpinner(opts)
	spinner.Complete()

	var buf bytes.Buffer
	spinner.Print(&buf, "✓")

	want := fmt.Sprintf("%s %s\r\n", "✓", initialMessage)
	assert.Equal(t, want, buf.String())
}

func TestPrintWithError(t *testing.T) {
	opts := initialOpts
	spinner := ysmrr.NewSpinner(opts)
	spinner.Error()

	var buf bytes.Buffer
	spinner.Print(&buf, "✗")

	want := fmt.Sprintf("%s %s\r\n", "✗", initialMessage)
	assert.Equal(t, want, buf.String())
}
