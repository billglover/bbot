package messsaging

import "testing"

func TestDefaultDestinationUnkown(t *testing.T) {
	e := Envelope{}

	if got, want := e.Destination.Type, UnknownDestination; got != want {
		t.Errorf("unexpected destination type: got %v, want %v", got, want)
	}
}
