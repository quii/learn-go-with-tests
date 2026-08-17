package poker

import (
	"bytes"
	"testing"
	"time"
)

func TestStdOutAlerter(t *testing.T) {
	t.Skip("this test is just for documentation, we don't want our build having an extra 5s added!")
	out := &bytes.Buffer{}
	alerter := StdOutAlerter(out)

	alerter.ScheduleAlertAt(5*time.Second, 100)

	time.Sleep(6 * time.Second)

	want := "Blind is now 100\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}
