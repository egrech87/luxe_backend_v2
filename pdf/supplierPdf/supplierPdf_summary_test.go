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

	// The area is the sum of actual width x actual height x quantity of each aperture,
	// the amount of apertures is the sum of the quantities. Child apertures are not part of the figures.
	main := aperture(2, 1000, 2000)
	main.ChildApertures = []*apertureEntity.ApertureEntity{aperture(5, 400, 500)}

	area, amount := summarise([]*apertureEntity.ApertureEntity{main, aperture(3, 900, 1200), aperture(1, 1500, 1500)})
	// 2 x 2.00 + 3 x 1.08 + 1 x 2.25
	if math.Abs(area-9.49) > 1e-9 {
		t.Fatalf("expected a total area of 9.49, got %v", area)
	}
	if amount != 6 {
		t.Fatalf("expected 6 apertures, got %d", amount)
	}

	area, amount = summarise(nil)
	if area != 0 || amount != 0 {
		t.Fatalf("expected no area and no apertures, got %v and %d", area, amount)
	}
}

// For a quote document the summary is calculated from width x height x quantity instead of the actual width and height,
// because Generate replaces them (withQuotedDimensions) before any section is built.
func TestSummariseQuoteDocumentUsesWidthAndHeight(t *testing.T) {
	aperture := func(quantity, width, height, actualWidth, actualHeight int) *apertureEntity.ApertureEntity {
		a := apertureEntity.New()
		a.Quantity, a.Width, a.Height, a.ActualWidth, a.ActualHeight = quantity, width, height, actualWidth, actualHeight
		return a
	}
	list := []*apertureEntity.ApertureEntity{
		aperture(1, 1000, 2000, 1100, 2100), // 2.00 quoted, 2.31 actual, quantity 1
		aperture(2, 900, 1000, 950, 1000),   // 0.90 quoted, 0.95 actual, quantity 2
	}

	area, amount := summarise(list)
	// 1 x 2.31 + 2 x 0.95
	if math.Abs(area-4.21) > 1e-9 || amount != 3 {
		t.Fatalf("order document: expected 4.21 and 3, got %v and %d", area, amount)
	}

	area, amount = summarise(withQuotedDimensions(list))
	// 1 x 2.00 + 2 x 0.90
	if math.Abs(area-3.80) > 1e-9 || amount != 3 {
		t.Fatalf("quote document: expected 3.80 and 3, got %v and %d", area, amount)
	}

	// The original apertures are not modified
	if list[0].ActualWidth != 1100 || list[0].ActualHeight != 2100 {
		t.Fatalf("the original aperture was modified: %dx%d", list[0].ActualWidth, list[0].ActualHeight)
	}
}
