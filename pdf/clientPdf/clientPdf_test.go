package clientPdf

import (
	"Quote/db/entity/entities/apertureEntity"
	"testing"
)

func TestSkuHeading(t *testing.T) {
	aperture := apertureEntity.New()
	aperture.ID = 1234
	aperture.UnitName, aperture.UnitNumber = "Window", 7

	aperture.Facade = "Living Room"
	if got, want := skuHeading(aperture), "SKU : Window 7 - Living Room - 1234"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}

	// No room (or only spaces): the room is left out
	for _, room := range []string{"", "   "} {
		aperture.Facade = room
		if got, want := skuHeading(aperture), "SKU : Window 7 - 1234"; got != want {
			t.Fatalf("room %q: got %q, want %q", room, got, want)
		}
	}
}
