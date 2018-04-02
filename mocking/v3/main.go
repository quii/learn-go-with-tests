package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Sleeper func(time.Duration)

const finalWord = "Go!"
const countdownStart = 5
const sleepDuration = 1 * time.Second

func Countdown(out io.Writer, sleep Sleeper) {
	for i := countdownStart; i > 0; i-- {
		sleep(sleepDuration)
		fmt.Fprintln(out, i)
	}

	sleep(sleepDuration)
	fmt.Fprint(out, finalWord)
}

func main() {
	Countdown(os.Stdout, time.Sleep)
}
