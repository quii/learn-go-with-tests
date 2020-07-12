package geometry

import "testing"

func TestPerimeter(t *testing.T) {
	rectangle := Rectangle{10.0, 10.0}
	got := Perimeter(rectangle)
	want := 40.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}

func TestArea(t *testing.T) {

	// checkArea := func(t *testing.T, shape Shape, got, want float64) {
	// 	t.Helper()
	// 	got = shape.Area()
	// 	if got != want {
	// 		t.Errorf("got %g want %g", got, want)
	// 	}

	// }

	areaTests := []struct {
		name  string
		shape Shape
		want  float64
	}{
		{name: "Rectangle", shape: Rectangle{12, 6}, want: 72.0},
		{name: "Circle", shape: Circle{10}, want: 314.1592653589793},
		{name: "Triangle", shape: Triangle{12, 6}, want: 36.0},
	}

	for _, test := range areaTests {
		t.Run(test.name, func(t *testing.T) {
			t.Helper()
			got := test.shape.Area()
			if got != test.want {
				t.Errorf("%#v got %g want %g", test.shape, got, test.want)
			}

		})
	}

	// t.Run("rectangle", func(t *testing.T) {
	// 	rectangle := Rectangle{6.0, 7.0}
	// 	got := rectangle.Area()
	// 	want := 42.0

	// 	checkArea(t, rectangle, got, want)
	// })

	// t.Run("circle", func(t *testing.T) {
	// 	circle := Circle{10}
	// 	got := circle.Area()
	// 	want := 314.1592653589793

	// 	checkArea(t, circle, got, want)
	// })

}
