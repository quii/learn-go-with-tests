package locks

import "testing"

func TestCounter(t *testing.T) {

	t.Run("incrementing the counter 3 times leaves it at 3", func(t *testing.T) {
		counter := Counter{}
		counter.Inc()
		counter.Inc()
		counter.Inc()

		assertCounter(t, counter, 3)
	})

	t.Run("it runs safely concurrently", func(t *testing.T) {
		wantedCount := 10
		counter := Counter{}

		for i:=0; i<wantedCount; i++ {
			go func() {
				counter.Inc()
			}()
		}

		assertCounter(t, counter, wantedCount)
	})

}

func assertCounter(t *testing.T, got Counter, want int)  {
	t.Helper()
	if got.Value() != want {
		t.Errorf("got %d, want %d", got.Value(), want)
	}
}
