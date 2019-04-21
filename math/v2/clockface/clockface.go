package clockface

import "time"

type Hands struct {
	Hour   Vector
	Minute Vector
	Second Vector
}

type Vector struct {
	X int
	Y int
}

func HandsAt(t time.Time) (hands Hands) {
	return
}
