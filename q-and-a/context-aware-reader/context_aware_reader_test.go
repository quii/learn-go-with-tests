package cancelreader

import (
	"context"
	"strings"
	"testing"
)

func TestContextAwareReader(t *testing.T) {
	t.Run("behaves like a normal reader", func(t *testing.T) {
		rdr := NewCancelableReader(context.Background(), strings.NewReader("123456"))
		got := make([]byte, 3)
		_, err := rdr.Read(got)

		if err != nil {
			t.Fatal(err)
		}

		assertBufferHas(t, got, "123")

		_, err = rdr.Read(got)

		if err != nil {
			t.Fatal(err)
		}

		assertBufferHas(t, got, "456")
	})

	t.Run("stops reading when canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		rdr := NewCancelableReader(ctx, strings.NewReader("123456"))
		got := make([]byte, 3)
		_, err := rdr.Read(got)

		if err != nil {
			t.Fatal(err)
		}

		assertBufferHas(t, got, "123")

		cancel()

		n, err := rdr.Read(got)

		if err == nil {
			t.Error("expected an error after cancellation but didnt get one")
		}

		if n > 0 {
			t.Errorf("expected 0 bytes to be read after cancellation but %d were read", n)
		}
	})
}

func assertBufferHas(t testing.TB, buf []byte, want string) {
	t.Helper()
	got := string(buf)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
