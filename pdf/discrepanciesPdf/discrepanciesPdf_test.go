package discrepanciesPdf

import (
	"Quote/db/entity/entities/apertureEntity"
	"Quote/db/entity/entities/colorEntity"
	"Quote/db/entity/entities/glassEntity"
	"Quote/db/entity/entities/profileEntity"
	"Quote/db/entity/entities/quoteEntity"
	"bytes"
	"math"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func aperture(id uint64, w, h, aw, ah int, price float64) *apertureEntity.ApertureEntity {
	a := apertureEntity.New()
	a.ID = id
	a.Width, a.Height, a.ActualWidth, a.ActualHeight = w, h, aw, ah
	a.DiscountedSellingPrice = price
	a.Profile = &profileEntity.ProfileEntity{FriendlyName: "Type"}
	a.Glass = &glassEntity.GlassEntity{FriendlyName: "Glass"}
	a.Color = &colorEntity.ColorEntity{FriendlyName: "Colour"}
	return a
}

func TestNewRow(t *testing.T) {
	row := NewRow(aperture(7, 1000, 2000, 1100, 2000, 100))
	if row.SKU != 7 || row.ApertureType != "Type" || row.Glass != "Glass" || row.Colour != "Colour" {
		t.Fatalf("unexpected row: %+v", row)
	}
	if math.Abs(row.QuotedPrice-200) > 1e-9 || math.Abs(row.ActualPrice-220) > 1e-9 || math.Abs(row.Discrepancy-20) > 1e-9 {
		t.Fatalf("unexpected prices: %+v", row)
	}

	// An aperture without a profile, glass or colour still gives a row
	bare := apertureEntity.New()
	if row := NewRow(bare); row.ApertureType != "" || row.Glass != "" || row.Colour != "" || row.Discrepancy != 0 {
		t.Fatalf("unexpected row for a bare aperture: %+v", row)
	}
}

func TestGeneratePaging(t *testing.T) {
	pages := func(count int) int {
		list := []*apertureEntity.ApertureEntity{}
		for i := 0; i < count; i++ {
			list = append(list, aperture(uint64(i+1), 1000, 2000, 1000, 2000, 100))
		}
		out, err := Generate(quoteEntity.New(), list)
		if err != nil {
			t.Fatal(err)
		}
		n, err := api.PageCount(bytes.NewReader(out), nil)
		if err != nil {
			t.Fatal(err)
		}
		return n
	}

	if n := pages(0); n != 1 {
		t.Fatalf("no apertures: expected 1 page, got %d", n)
	}
	if n := pages(10); n != 1 {
		t.Fatalf("10 apertures: expected 1 page, got %d", n)
	}
	if n := pages(100); n < 2 {
		t.Fatalf("100 apertures: expected several pages, got %d", n)
	}
}

func TestReadableMoney(t *testing.T) {
	for amount, want := range map[float64]string{1234.5: "\x801,234.50", -300.714: "-\x80300.71", -0.001: "\x800.00", 0: "\x800.00", -1234567.891: "-\x801,234,567.89"} {
		if got := readableMoney(amount); got != want {
			t.Fatalf("readableMoney(%v) = %q, want %q", amount, got, want)
		}
	}
}
