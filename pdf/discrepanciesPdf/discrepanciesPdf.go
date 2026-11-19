package discrepanciesPdf

import (
	"Quote/db/entity/entities/apertureEntity"
	"Quote/db/entity/entities/quoteEntity"
	"Quote/utils/formatters"
	"bytes"
	"fmt"
	"math"
	"time"

	"codeberg.org/go-pdf/fpdf"
)

// Row holds the values shown for one aperture of the discrepancies table.
type Row struct {
	SKU          uint64
	FriendlyName string
	Unit         string
	ApertureType string
	Glass        string
	Colour       string
	Width        int
	Height       int
	ActualWidth  int
	ActualHeight int
	ActualPrice  float64
	QuotedPrice  float64
	Discrepancy  float64
}

// NewRow builds the row of an aperture. The prices are an area in square metres (width x height, in mm) multiplied by the
// aperture's discounted selling price: the quoted price uses the quoted width and height, the actual price the actual ones.
// The discrepancy is the actual price minus the quoted price.
func NewRow(aperture *apertureEntity.ApertureEntity) Row {
	row := Row{
		SKU:          aperture.GetID(),
		FriendlyName: aperture.GetFriendlyName(),
		Unit:         aperture.GetUnitFullName(),
		Width:        aperture.GetWidth(),
		Height:       aperture.GetHeight(),
		ActualWidth:  aperture.GetActualWidth(),
		ActualHeight: aperture.GetActualHeight(),
	}
	if aperture.Profile != nil {
		row.ApertureType = aperture.Profile.FriendlyName
	}
	if aperture.Glass != nil {
		row.Glass = aperture.Glass.FriendlyName
	}
	if aperture.Color != nil {
		row.Colour = aperture.Color.FriendlyName
	}

	price := aperture.GetDiscountedSellingPrice()
	row.QuotedPrice = float64(row.Width) * float64(row.Height) / 1000000.0 * price
	row.ActualPrice = float64(row.ActualWidth) * float64(row.ActualHeight) / 1000000.0 * price
	row.Discrepancy = row.ActualPrice - row.QuotedPrice

	return row
}

// Generate creates the discrepancies PDF: a table with a row for each of the given (main) apertures,
// followed by a row with the total actual price, total quoted price and total discrepancy.
func Generate(quote *quoteEntity.QuoteEntity, apertureList []*apertureEntity.ApertureEntity) ([]byte, error) {
	const (
		margin     = 10.0
		lineHt     = 4.2
		cellPad    = 1.0
		headerHt   = 13.0
		bottomLine = 190.0 // A4 landscape is 210 mm high; rows that would pass this line start a new page
	)

	type column struct {
		Title string
		Wd    float64
		Align string // "L" for text, "R" for numbers
	}
	columns := []column{
		{"SKU", 12, "L"},
		{"Friendly Name", 30, "L"},
		{"Unit", 18, "L"},
		{"Aperture Type", 28, "L"},
		{"Glass", 28, "L"},
		{"Colour", 22, "L"},
		{"Width (mm)", 16, "R"},
		{"Height (mm)", 16, "R"},
		{"Actual Width (mm)", 19, "R"},
		{"Actual Height (mm)", 19, "R"},
		{"Actual Price", 24, "R"},
		{"Quoted Price", 24, "R"},
		{"Discrepancy", 21, "R"},
	}

	pdf := fpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(margin, margin, margin)
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()

	y := margin

	// Title
	pdf.SetFont("Helvetica", "B", 14)
	pdf.SetTextColor(15, 36, 62)
	pdf.SetXY(margin, y)
	pdf.CellFormat(150, 6, "Discrepancies", "", 0, "TL", false, 0, "")

	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(77, 77, 79)
	pdf.SetXY(margin+150, y)
	pdf.CellFormat(297-2*margin-150, 6, fmt.Sprintf("Quote : %d     Date : %s", quote.GetNumber(), time.Now().Format("02.01.2006")), "", 0, "TR", false, 0, "")
	y += 10

	drawHeader := func() {
		x := margin
		pdf.SetFont("Helvetica", "B", 8)
		pdf.SetFillColor(155, 194, 230)
		pdf.SetDrawColor(77, 77, 79)
		pdf.SetTextColor(15, 36, 62)
		pdf.SetLineWidth(0.3)

		for _, col := range columns {
			pdf.Rect(x, y, col.Wd, headerHt, "FD")

			lines := pdf.SplitLines([]byte(col.Title), col.Wd-2*cellPad)
			textY := y + (headerHt-float64(len(lines))*lineHt)/2
			for i, line := range lines {
				pdf.SetXY(x+cellPad, textY+float64(i)*lineHt)
				pdf.CellFormat(col.Wd-2*cellPad, lineHt, string(line), "", 0, "CM", false, 0, "")
			}
			x += col.Wd
		}
		y += headerHt
	}
	drawHeader()

	var totalActualPrice, totalQuotedPrice, totalDiscrepancy float64

	for _, aperture := range apertureList {
		row := NewRow(aperture)
		totalActualPrice += row.ActualPrice
		totalQuotedPrice += row.QuotedPrice
		totalDiscrepancy += row.Discrepancy
		texts := []string{
			fmt.Sprintf("%d", row.SKU),
			row.FriendlyName,
			row.Unit,
			row.ApertureType,
			row.Glass,
			row.Colour,
			formatters.ToReadableWholeNumber(row.Width),
			formatters.ToReadableWholeNumber(row.Height),
			formatters.ToReadableWholeNumber(row.ActualWidth),
			formatters.ToReadableWholeNumber(row.ActualHeight),
			readableMoney(row.ActualPrice),
			readableMoney(row.QuotedPrice),
			readableMoney(row.Discrepancy),
		}

		// The row is as high as its tallest cell, since the text cells may wrap
		pdf.SetFont("Helvetica", "", 8)
		cellLines := make([][][]byte, len(columns))
		rowHt := lineHt + 2*cellPad
		for i, col := range columns {
			cellLines[i] = pdf.SplitLines([]byte(texts[i]), col.Wd-2*cellPad)
			if ht := float64(len(cellLines[i]))*lineHt + 2*cellPad; ht > rowHt {
				rowHt = ht
			}
		}

		if y+rowHt > bottomLine {
			pdf.AddPage()
			y = margin
			drawHeader()
			pdf.SetFont("Helvetica", "", 8)
		}

		x := margin
		pdf.SetDrawColor(160, 160, 160)
		pdf.SetLineWidth(0.2)
		pdf.SetTextColor(0, 0, 0)
		for i, col := range columns {
			pdf.Rect(x, y, col.Wd, rowHt, "D")
			for j, line := range cellLines[i] {
				pdf.SetXY(x+cellPad, y+cellPad+float64(j)*lineHt)
				pdf.CellFormat(col.Wd-2*cellPad, lineHt, string(line), "", 0, col.Align+"M", false, 0, "")
			}
			x += col.Wd
		}
		y += rowHt
	}

	// The totals row, at the bottom of the table
	{
		totalHt := lineHt + 2*cellPad
		if y+totalHt > bottomLine {
			pdf.AddPage()
			y = margin
			drawHeader()
		}

		pdf.SetFont("Helvetica", "B", 8)
		pdf.SetFillColor(225, 235, 247)
		pdf.SetDrawColor(77, 77, 79)
		pdf.SetTextColor(15, 36, 62)
		pdf.SetLineWidth(0.3)

		// The label spans all the columns before the three price columns
		const priceColumns = 3
		labelColumns := columns[:len(columns)-priceColumns]
		labelWd := 0.0
		for _, col := range labelColumns {
			labelWd += col.Wd
		}
		x := margin
		pdf.Rect(x, y, labelWd, totalHt, "FD")
		pdf.SetXY(x+cellPad, y+cellPad)
		pdf.CellFormat(labelWd-2*cellPad, lineHt, "Total", "", 0, "RM", false, 0, "")
		x += labelWd

		totals := []float64{totalActualPrice, totalQuotedPrice, totalDiscrepancy}
		for i, total := range totals {
			col := columns[len(labelColumns)+i]
			pdf.Rect(x, y, col.Wd, totalHt, "FD")
			pdf.SetXY(x+cellPad, y+cellPad)
			pdf.CellFormat(col.Wd-2*cellPad, lineHt, readableMoney(total), "", 0, "RM", false, 0, "")
			x += col.Wd
		}
		y += totalHt
	}

	var buff bytes.Buffer
	if err := pdf.Output(&buff); err != nil {
		return []byte{}, err
	}

	return buff.Bytes(), nil
}

// euro is the Euro sign in the Windows-1252 encoding used by the PDF's core fonts.
const euro = "\x80"

// readableMoney formats an amount in Euro with thousands separators and two decimals, keeping the sign in front: -€1,234.50.
func readableMoney(amount float64) string {
	if math.Abs(amount) < 0.005 {
		amount = 0 // avoid showing "-€0.00"
	}
	if amount < 0 {
		return "-" + euro + formatters.ToReadableNumber_2dp(-amount)
	}
	return euro + formatters.ToReadableNumber_2dp(amount)
}
