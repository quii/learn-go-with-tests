package qanda

import (
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
	"testing"
)

type Payload struct {
	Message string `xml:"message"`
}

func KeithCommand(data io.Reader) string {
	var payload Payload
	xml.NewDecoder(data).Decode(&payload)
	return strings.ToUpper(payload.Message)
}

func getXMLFromCommand() io.Reader {

	cmd := exec.Command("cat", "msg.xml")
	out, _ := cmd.StdoutPipe()

	cmd.Start()
	data, _ := ioutil.ReadAll(out)
	cmd.Wait()

	return bytes.NewReader(data)
}

func TestKeithCommandIntegration(t *testing.T) {
	got := KeithCommand(getXMLFromCommand())
	want := "HAPPY NEW YEAR!"

	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}

func TestKeithCommand(t *testing.T) {

	t.Run("an example", func(t *testing.T) {

		input := strings.NewReader(`
<payload>
    <message>Cats are the best animal</message>
</payload>`)

		got := KeithCommand(input)
		want := "CATS ARE THE BEST ANIMAL"

		if got != want {
			t.Errorf("got '%s', want '%s'", got, want)
		}
	})
}
