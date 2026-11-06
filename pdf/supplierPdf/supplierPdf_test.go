package supplierPdf

import (
	"Quote/db/entity/entities/apertureEntity"
	"testing"
)

func TestSplitIntoUnits(t *testing.T) {
	child := apertureEntity.New()
	child.Quantity = 2
	child.ActualWidth = []int64{400, 410}
	child.ActualHeight = []int64{500, 500}

	main := apertureEntity.New()
	main.Quantity = 2
	main.ActualWidth = []int64{900, 905}
	main.ActualHeight = []int64{1200, 1200}
	main.ChildApertures = []*apertureEntity.ApertureEntity{child}

	units := splitIntoUnits([]*apertureEntity.ApertureEntity{main})
	if len(units) != 2 {
		t.Fatalf("expected 2 units, got %d", len(units))
	}
	for i, u := range units {
		if u.Quantity != 1 || len(u.ActualWidth) != 1 || u.ActualWidth[0] != main.ActualWidth[i] {
			t.Fatalf("unit %d wrong: %+v", i, u)
		}
		if len(u.ChildApertures) != 1 || u.ChildApertures[0].ActualWidth[0] != child.ActualWidth[i] || u.ChildApertures[0].Quantity != 1 {
			t.Fatalf("unit %d child wrong", i)
		}
	}
	if main.Quantity != 2 || len(main.ActualWidth) != 2 || len(main.ChildApertures) != 1 {
		t.Fatal("input was modified")
	}
}
