package statuswriter

import (
	"time"

	"github.com/jedib0t/go-pretty/v6/progress"
)

func DisableColors() {
	defaultColorStyle = progress.StyleColorsDefault
}

type StatusWriter struct {
	Writer progress.Writer
}

func New() *StatusWriter {
	return newStatusWriterWithDefaults()
}

var defaultColorStyle = progress.StyleColorsExample

const (
	defaultSleepDuration   = time.Millisecond * 100
	defaultUpdateFrequency = time.Millisecond * 100
)

func newStatusWriterWithDefaults() *StatusWriter {
	t := &StatusWriter{Writer: progress.NewWriter()}
	t.Writer.ShowOverallTracker(false)
	t.Writer.SetStyle(progress.StyleDefault)
	t.Writer.SetTrackerPosition(progress.PositionRight)
	t.Writer.ShowTime(true)
	t.Writer.ShowValue(false)
	t.Writer.ShowPercentage(false)
	t.Writer.ShowTracker(false)
	t.Writer.SetUpdateFrequency(defaultUpdateFrequency)
	t.Writer.Style().Colors = defaultColorStyle
	t.Writer.Style().Options.DoneString = ""
	t.Writer.Style().Options.Separator = ""
	t.Writer.Style().Options.TimeInProgressPrecision = time.Millisecond
	t.Writer.Style().Options.TimeDonePrecision = time.Millisecond
	return t
}

func (t StatusWriter) Run() {
	go t.Writer.Render()
}

func (t StatusWriter) Stop() {
	t.Writer.Stop()
}

func (t StatusWriter) Wait() {
	time.Sleep(defaultSleepDuration)
	for t.Writer.IsRenderInProgress() {
		if t.Writer.LengthActive() == 0 {
			t.Writer.Stop()
		}
		time.Sleep(defaultSleepDuration)
	}
}

type Sink struct {
	tracker *progress.Tracker
}

func (t StatusWriter) NewSink() *Sink {
	tracker := &progress.Tracker{}
	t.Writer.AppendTracker(tracker)
	return &Sink{tracker: tracker}
}

func (s Sink) Write(status string) {
	s.tracker.Message = status
}

func (s Sink) Close() {
	s.tracker.MarkAsDone()
}
