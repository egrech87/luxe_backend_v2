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

func TestSplitIntoUnitsGroupsIdentical(t *testing.T) {
	main := apertureEntity.New()
	main.Quantity = 4
	main.ActualWidth = []int64{900, 905, 900, 900}
	main.ActualHeight = []int64{1200, 1200, 1200, 1200}

	units := splitIntoUnits([]*apertureEntity.ApertureEntity{main})
	if len(units) != 2 {
		t.Fatalf("expected 2 items, got %d", len(units))
	}
	if units[0].Quantity != 3 || len(units[0].ActualWidth) != 3 || units[0].ActualWidth[0] != 900 {
		t.Fatalf("first item wrong: %+v", units[0])
	}
	if units[1].Quantity != 1 || units[1].ActualWidth[0] != 905 {
		t.Fatalf("second item wrong: %+v", units[1])
	}
	if units[0].GetActualAreaSqm() != 3*0.9*1.2 {
		t.Fatalf("area wrong: %v", units[0].GetActualAreaSqm())
	}
}
