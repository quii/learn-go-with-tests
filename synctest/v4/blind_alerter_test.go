package poker

import (
	"bytes"
	"testing"
	"testing/synctest"
	"time"
)

func TestStdOutAlerter(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		out := &bytes.Buffer{}
		alerter := StdOutAlerter(out)

		alerter.ScheduleAlertAt(5*time.Second, 100)

		synctest.Wait()
		if out.String() != "" {
			t.Errorf("did not expect anything to be printed yet, got %q", out.String())
		}

		time.Sleep(5 * time.Second)
		synctest.Wait()

		want := "Blind is now 100\n"
		if out.String() != want {
			t.Errorf("got %q, want %q", out.String(), want)
		}
	})
}
