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

func HandsAt(t time.Time) Hands {
	return Hands{
		Hour:   Vector{X: 0, Y: 150},
		Minute: Vector{X: 0, Y: 150},
		Second: Vector{X: 0, Y: 150},
	}
}
