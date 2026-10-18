package poker

import (
	"fmt"
	"time"
)

// BlindAlerter schedules alerts for blind amounts.
type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int)
}

// BlindAlerterFunc allows you to implement BlindAlerter with a function.
type BlindAlerterFunc func(duration time.Duration, amount int)

// ScheduleAlertAt is BlindAlerterFunc's implementation of BlindAlerter.
func (a BlindAlerterFunc) ScheduleAlertAt(duration time.Duration, amount int) {
	a(duration, amount)
}

// NewAlerter returns a BlindAlerterFunc and a channel that receives the
// formatted alert message once each scheduled duration elapses.
func NewAlerter() (BlindAlerterFunc, <-chan string) {
	alerts := make(chan string)

	scheduleAlertAt := func(duration time.Duration, amount int) {
		time.AfterFunc(duration, func() {
			alerts <- fmt.Sprintf("Blind is now %d", amount)
		})
	}

	return scheduleAlertAt, alerts
}
