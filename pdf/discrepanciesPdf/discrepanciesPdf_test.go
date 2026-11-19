package discrepanciesPdf

import (
	"Quote/db/entity/entities/apertureEntity"
	"Quote/db/entity/entities/colorEntity"
	"Quote/db/entity/entities/glassEntity"
	"Quote/db/entity/entities/profileEntity"
	"Quote/db/entity/entities/quoteEntity"
	"bytes"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func aperture(id uint64, w, h, aw, ah int, price float64) *apertureEntity.ApertureEntity {
	a := apertureEntity.New()
	a.ID = id
	a.Width, a.Height, a.ActualWidth, a.ActualHeight = w, h, aw, ah
	a.DiscountedSellingPrice = price
	a.FriendlyName = "Name"
	a.Profile = &profileEntity.ProfileEntity{FriendlyName: "Type"}
	a.Glass = &glassEntity.GlassEntity{FriendlyName: "Glass"}
	a.Color = &colorEntity.ColorEntity{FriendlyName: "Colour"}
	return a
}

func TestNewRow(t *testing.T) {
	row := NewRow(aperture(7, 1000, 2000, 1100, 2000, 100))
	if row.SKU != 7 || row.FriendlyName != "Name" || row.ApertureType != "Type" || row.Glass != "Glass" || row.Colour != "Colour" {
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

func TestGenerateShowsTotals(t *testing.T) {
	if _, err := exec.LookPath("pdftotext"); err != nil {
		t.Skip("pdftotext is not available")
	}

	out, err := Generate(quoteEntity.New(), []*apertureEntity.ApertureEntity{
		aperture(1, 1000, 2000, 1100, 2000, 100),    // actual 220.00, quoted 200.00
		aperture(3, 1200, 2100, 1150, 2050, 1850.5), // actual 4,362.55, quoted 4,663.26
		aperture(4, 900, 900, 900, 900, 333.33),     // actual 270.00, quoted 270.00
	})
	if err != nil {
		t.Fatal(err)
	}

	file := filepath.Join(t.TempDir(), "discrepancies.pdf")
	if err := os.WriteFile(file, out, 0644); err != nil {
		t.Fatal(err)
	}
	text, err := exec.Command("pdftotext", "-layout", file, "-").Output()
	if err != nil {
		t.Fatal(err)
	}

	var totalLine string
	for _, line := range strings.Split(string(text), "\n") {
		if strings.Contains(line, "Total") {
			totalLine = line
		}
	}
	// The euro sign is not part of the extracted text encoding, so only the amounts are compared
	for _, want := range []string{"4,852.55", "5,133.26", "-", "280.71"} {
		if !strings.Contains(totalLine, want) {
			t.Fatalf("totals row %q does not contain %q", totalLine, want)
		}
	}
}
