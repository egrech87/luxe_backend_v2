package supplierPdf

import (
	"Quote/db/entity/entities/apertureEntity"
	"math"
	"testing"
)

func TestSummarise(t *testing.T) {
	aperture := func(quantity, actualWidth, actualHeight int) *apertureEntity.ApertureEntity {
		a := apertureEntity.New()
		a.Quantity, a.ActualWidth, a.ActualHeight = quantity, actualWidth, actualHeight
		return a
	}

	// The area is the sum of actual width x actual height of each aperture (not multiplied by the quantity),
	// the amount of apertures is the sum of the quantities. Child apertures are not part of the figures.
	main := aperture(2, 1000, 2000)
	main.ChildApertures = []*apertureEntity.ApertureEntity{aperture(5, 400, 500)}

	area, amount := summarise([]*apertureEntity.ApertureEntity{main, aperture(3, 900, 1200), aperture(1, 1500, 1500)})
	if math.Abs(area-5.33) > 1e-9 {
		t.Fatalf("expected a total area of 5.33, got %v", area)
	}
	if amount != 6 {
		t.Fatalf("expected 6 apertures, got %d", amount)
	}

	area, amount = summarise(nil)
	if area != 0 || amount != 0 {
		t.Fatalf("expected no area and no apertures, got %v and %d", area, amount)
	}
}

// For a quote document the summary is calculated from the width and height instead of the actual width and height,
// because Generate replaces them (withQuotedDimensions) before any section is built.
func TestSummariseQuoteDocumentUsesWidthAndHeight(t *testing.T) {
	aperture := func(quantity, width, height, actualWidth, actualHeight int) *apertureEntity.ApertureEntity {
		a := apertureEntity.New()
		a.Quantity, a.Width, a.Height, a.ActualWidth, a.ActualHeight = quantity, width, height, actualWidth, actualHeight
		return a
	}
	list := []*apertureEntity.ApertureEntity{
		aperture(1, 1000, 2000, 1100, 2100), // 2.00 quoted, 2.31 actual
		aperture(2, 900, 1000, 950, 1000),   // 0.90 quoted, 0.95 actual
	}

	area, amount := summarise(list)
	if math.Abs(area-3.26) > 1e-9 || amount != 3 {
		t.Fatalf("order document: expected 3.26 and 3, got %v and %d", area, amount)
	}

	area, amount = summarise(withQuotedDimensions(list))
	if math.Abs(area-2.90) > 1e-9 || amount != 3 {
		t.Fatalf("quote document: expected 2.90 and 3, got %v and %d", area, amount)
	}

	// The original apertures are not modified
	if list[0].ActualWidth != 1100 || list[0].ActualHeight != 2100 {
		t.Fatalf("the original aperture was modified: %dx%d", list[0].ActualWidth, list[0].ActualHeight)
	}
}
