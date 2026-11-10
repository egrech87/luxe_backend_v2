package supplierPdf

import (
	"Quote/db/entity/entities/apertureEntity"
	"Quote/db/entity/entities/colorEntity"
	"Quote/db/entity/entities/quoteEntity"
	"Quote/model/enums/fileTypes"
	"Quote/utils/formatters"
	"Quote/utils/pathUtils"
	"Quote/utils/pdfUtils"
	"bytes"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"

	"codeberg.org/go-pdf/fpdf"
	"github.com/pkg/errors"
)

// Generate builds the supplier PDF. If termsAndConditions is not empty, that PDF is merged at the very end of the report, after the datasheets.
func Generate(quote *quoteEntity.QuoteEntity, mainApertureList []*apertureEntity.ApertureEntity, datasheetList [][]byte, termsAndConditions []byte) ([]byte, error) {
	if err := validate(quote, mainApertureList); err != nil {
		return []byte{}, err
	}

	// var color *colorEntity.ColorEntity
	// if len(mainApertureList) > 0 {
	// 	color = mainApertureList[0].Color
	// }

	sections := [][]byte{}

	// section1, err := generateSection1(quote, color)
	// if err != nil {
	// 	return []byte{}, err
	// }
	// sections = append(sections, section1)

	section2, err := generateSection2(mainApertureList)
	if err != nil {
		return []byte{}, err
	}
	sections = append(sections, section2)

	section3, err := generateSection3(quote, mainApertureList)
	if err != nil {
		return []byte{}, err
	}
	sections = append(sections, section3)

	sections = append(sections, datasheetList...)

	if len(termsAndConditions) > 0 {
		sections = append(sections, termsAndConditions)
	}

	fullPdf, err := pdfUtils.MergePdfBytes(sections...)
	if err != nil {
		return []byte{}, err
	}

	// Number every page of the final document
	fullPdf, err = pdfUtils.AddPageNumbers(fullPdf)
	if err != nil {
		return []byte{}, err
	}

	return fullPdf, err
}

func validate(quote *quoteEntity.QuoteEntity, mainApertureList []*apertureEntity.ApertureEntity) error {
	errorMessage := "Cannot generate quote since"
	if quote.Client == nil {
		return errors.Errorf("%s the quote has no client", errorMessage)
	}

	if len(mainApertureList) == 0 {
		return errors.Errorf("%s the quote has no apertures", errorMessage)
	}

	for _, parentAperture := range mainApertureList {
		if parentAperture.Profile == nil {
			return errors.Errorf("%s the aperture %s has no aperture type", errorMessage, parentAperture.FriendlyName)
		}
		if parentAperture.Glass == nil {
			return errors.Errorf("%s the aperture %s has no glass", errorMessage, parentAperture.FriendlyName)
		}
		if parentAperture.Color == nil {
			return errors.Errorf("%s the aperture %s has no color", errorMessage, parentAperture.FriendlyName)
		}

		for _, childAperture := range parentAperture.ChildApertures {
			if childAperture.Profile == nil {
				return errors.Errorf("%s the aperture %s has no aperture type", errorMessage, childAperture.FriendlyName)
			}
			if childAperture.Glass == nil {
				return errors.Errorf("%s the aperture %s has no glass", errorMessage, childAperture.FriendlyName)
			}
			if childAperture.Color == nil {
				return errors.Errorf("%s the aperture %s has no color", errorMessage, childAperture.FriendlyName)
			}
		}
	}

	return nil
}

// This is page 1
func generateSection1(quote *quoteEntity.QuoteEntity, color *colorEntity.ColorEntity) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(5, 5, 5)
	pdf.AddPage()
	pdf.SetFont("Helvetica", "", 12)

	drawWatermark := func() {
		pdf.SetFillColor(243, 245, 247)
		pdf.Rect(80, 0, 41.7, 191, "F")
		pdf.Rect(40.5, 112, 169.5, 41.7, "F")
	}

	fillTopLeftSection := func() {
		const (
			x        float64 = 10.4
			imgY     float64 = 5.3
			imgWd    float64 = 32.7
			imgHt    float64 = 19
			lnHt     float64 = 1.9
			line1Y   float64 = imgY + imgHt + 3.1
			line2Y   float64 = line1Y + lnHt + 3.1
			line3Y   float64 = line2Y + lnHt + 3.1
			line4Y   float64 = line3Y + lnHt + 3.1
			vatLineY float64 = line4Y + lnHt + 12.5
			clientY  float64 = vatLineY + lnHt + 7.7
		)

		// Draw the logo
		var opt fpdf.ImageOptions
		opt.ImageType = "png"
		pdf.ImageOptions(pathUtils.ToImagePath("luxe-logo.PNG"), x, imgY, imgWd, imgHt, false, opt, 0, "")

		pdf.SetFontSize(10)

		// Line 1
		pdf.SetFontStyle("B")
		pdf.SetTextColor(15, 36, 62)
		pdf.SetXY(x, line1Y)
		pdf.CellFormat(100, 2, "LUXE Apertures Ltd.", "", 0, "TL", false, 0, "")

		// Line 2
		pdf.SetFontStyle("")
		pdf.SetTextColor(140, 140, 140)
		pdf.SetXY(x, line2Y)
		pdf.CellFormat(100, 2, "75, Triq tal-Gebel", "", 0, "TL", false, 0, "")

		// Line 3
		pdf.SetXY(x, line3Y)
		pdf.CellFormat(100, 2, "Birzebbuga", "", 0, "TL", false, 0, "")

		// Line 4
		pdf.SetXY(x, line4Y)
		pdf.CellFormat(100, 2, "BBG2775", "", 0, "TL", false, 0, "")

		// Line 4
		pdf.SetXY(x, vatLineY)
		pdf.CellFormat(100, 2, "VAT Reg No:  MT29023304", "", 0, "TL", false, 0, "")

		// Client section
		clientFullname := quote.Client.GetFullName()
		addressLines := strings.Split(quote.Client.Address, "\n")
		if len(clientFullname) > 0 {
			addressLines = slices.Insert(addressLines, 0, clientFullname)
		}

		for i, ln := range addressLines {
			pdf.SetXY(x, clientY+float64(i)*(lnHt+3.1))
			pdf.CellFormat(100, 2, ln, "", 0, "TL", false, 0, "")
		}
	}

	fillTopRightSection := func() {
		const line1Y float64 = 32.2

		pdf.SetFontSize(10)
		pdf.SetTextColor(140, 140, 140)

		// Line 1 - Quotation word
		pdf.SetFontStyle("")
		pdf.SetXY(122.7, line1Y)
		pdf.CellFormat(100, 2, "Quotation", "", 0, "TL", false, 0, "")

		// Line 1 - Page word
		pdf.SetFontStyle("B")
		pdf.SetXY(176, line1Y)
		pdf.CellFormat(100, 2, "Page", "", 0, "TL", false, 0, "")

		// Line 1 - The no. 1
		pdf.SetFontStyle("")
		pdf.SetXY(186.3, line1Y)
		pdf.CellFormat(100, 2, "1", "", 0, "TL", false, 0, "")

		// Quotation no. Line
		pdf.SetFontStyle("")
		pdf.SetXY(122.7, 66.4)
		pdf.CellFormat(39, 2, "NUMBER", "", 0, "TR", false, 0, "")
		pdf.SetXY(164.5, 66.4)
		pdf.CellFormat(5, 2, strconv.Itoa(quote.Number), "", 0, "TL", false, 0, "")

		// Date Line
		now := time.Now()
		pdf.SetFontStyle("")
		pdf.SetXY(122.7, 74)
		pdf.CellFormat(39, 2, "DATE", "", 0, "TR", false, 0, "")
		pdf.SetXY(164.5, 74)
		pdf.CellFormat(5, 2, now.Format("02.01.2006"), "", 0, "TL", false, 0, "")

		// Expiry Date Line
		expiry := now.AddDate(0, 1, 0)
		pdf.SetFontStyle("")
		pdf.SetXY(122.7, 81.6)
		pdf.CellFormat(39, 2, "EXPIRY DATE", "", 0, "TR", false, 0, "")
		pdf.SetXY(164.5, 81.6)
		pdf.CellFormat(5, 2, expiry.Format("02.01.2006"), "", 0, "TL", false, 0, "")
	}

	fillProductDescription := func() {
		const (
			titleY         float64 = 107
			quoteAddresssY float64 = titleY + 11.4
			colorY         float64 = quoteAddresssY + 4.5
			netsY          float64 = colorY + 4.5
		)

		// Title line
		pdf.SetFontStyle("B")
		pdf.SetXY(32.2, titleY)
		pdf.CellFormat(100, 2, "Product Description", "", 0, "TL", false, 0, "")

		// Quote address line
		quoteAddress := ""
		{
			addressParts := []string{}
			if quote.Street != "" {
				addressParts = append(addressParts, quote.Street)
			}
			if quote.Locality != "" {
				addressParts = append(addressParts, quote.Locality)
			}
			if quote.Country != "" {
				addressParts = append(addressParts, quote.Country)
			}
			quoteAddress = strings.Join(addressParts, ", ")
		}
		pdf.SetFontStyle("")
		pdf.SetFontSize(8)
		pdf.SetXY(32.2, quoteAddresssY)
		pdf.CellFormat(100, 2, quoteAddress, "", 0, "TL", false, 0, "")

		// Color line
		colorName := "N/A"
		if color != nil {
			colorName = color.FriendlyName
		}
		pdf.SetXY(32.2, colorY)
		pdf.CellFormat(100, 2, "Colour : "+colorName, "", 0, "TL", false, 0, "")

		// Nets line
		incMsg := "Included"
		if !quote.HasNets() {
			incMsg = "Not Included"
		}
		pdf.SetXY(32.2, netsY)
		pdf.CellFormat(100, 2, "Nets : "+incMsg, "", 0, "TL", false, 0, "")
	}

	fillBottomLeftSection := func() {
		const (
			col1X      float64 = 10
			col2X      float64 = 56
			listIndent float64 = 2.8
			line1Y             = 249.2
			line2Y             = line1Y + 3.4
			line3Y             = line2Y + 3.4
			line4Y             = line3Y + 3.4
			line5Y             = line4Y + 3.4
			line6Y             = line5Y + 3.4
			line7Y             = line6Y + 3.4
		)

		pdf.SetTextColor(165, 165, 165)
		pdf.SetFontStyle("")
		pdf.SetFontSize(8)

		// Line 1
		pdf.SetXY(col1X, line1Y)
		pdf.CellFormat(100, 2, "PVC Guarantee:", "", 0, "TL", false, 0, "")

		pdf.SetXY(col2X, line1Y)
		pdf.CellFormat(100, 2, "Payment Terms:", "", 0, "TL", false, 0, "")

		// Line 2
		pdf.SetXY(col1X+listIndent, line2Y)
		pdf.CellFormat(100, 2, "- White PVC - 10 Years", "", 0, "TL", false, 0, "")

		pdf.SetXY(col2X+listIndent, line2Y)
		pdf.CellFormat(100, 2, "- 50% Deposit", "", 0, "TL", false, 0, "")

		// Line 3
		pdf.SetXY(col1X+listIndent, line3Y)
		pdf.CellFormat(100, 2, "- Coloured PVC - 5 Years", "", 0, "TL", false, 0, "")

		pdf.SetXY(col2X+listIndent, line3Y)
		pdf.CellFormat(100, 2, "- 25% Material on Site", "", 0, "TL", false, 0, "")

		// Line 4
		pdf.SetXY(col2X+listIndent, line4Y)
		pdf.CellFormat(100, 2, "- 25% Measured", "", 0, "TL", false, 0, "")

		// Line 5
		pdf.SetXY(col1X, line5Y)
		pdf.CellFormat(100, 2, "Delivery: 8 weeks from deposit (subject that final measurements", "", 0, "TL", false, 0, "")

		// Line 6
		pdf.SetXY(col1X, line6Y)
		pdf.CellFormat(100, 2, "are available)", "", 0, "TL", false, 0, "")

		// Line 7
		pdf.SetXY(col1X, line7Y)
		pdf.CellFormat(100, 2, "Permits and Lifting Services are excluded", "", 0, "TL", false, 0, "")
	}

	drawWatermark()
	fillTopLeftSection()
	fillTopRightSection()
	fillProductDescription()
	fillBottomLeftSection()

	var buff bytes.Buffer
	err := pdf.Output(&buff)
	if err != nil {
		return []byte{}, err
	}

	return buff.Bytes(), nil
}

func generateSection2(mainApertureList []*apertureEntity.ApertureEntity) ([]byte, error) {
	pdf := fpdf.New("L", "mm", "Letter", "")
	pdf.SetMargins(5, 5, 5)
	pdf.AddPage()
	pdf.SetFont("Helvetica", "", 12)

	type ColInfo struct {
		X  float64
		Wd float64
	}

	col1 := ColInfo{
		X:  6.5,
		Wd: 12.3,
	}
	col2 := ColInfo{
		X:  18.8,
		Wd: 13.7,
	}
	col3 := ColInfo{
		X:  32.5,
		Wd: 15.2,
	}
	col4 := ColInfo{
		X:  47.7,
		Wd: 49.3,
	}
	col5 := ColInfo{
		X:  97,
		Wd: 26.5,
	}
	col6 := ColInfo{
		X:  123.5,
		Wd: 15.9,
	}
	col7 := ColInfo{
		X:  139.4,
		Wd: 13.9,
	}
	col8 := ColInfo{
		X:  153.3,
		Wd: 12.4,
	}
	col9 := ColInfo{
		X:  165.7,
		Wd: 22.7,
	}
	col10 := ColInfo{
		X:  188.4,
		Wd: 11.2,
	}
	col11 := ColInfo{
		X:  199.6,
		Wd: 16.2,
	}
	col12 := ColInfo{
		X:  215.8,
		Wd: 16.2,
	}
	col13 := ColInfo{
		X:  232,
		Wd: 13.4,
	}
	col14 := ColInfo{
		X:  245.4,
		Wd: 24.1,
	}
	colArr := [14]ColInfo{
		col1,
		col2,
		col3,
		col4,
		col5,
		col6,
		col7,
		col8,
		col9,
		col10,
		col11,
		col12,
		col13,
		col14,
	}

	var y float64
	pageContainsHeader := true

	drawTopHorizontalLine := func() {
		pdf.SetLineWidth(0.4)
		pdf.Line(col1.X, 19.3, 269.5, 19.3)
	}

	drawVerticalLines := func(endY float64) {
		y1 := 19.3
		y2 := 19.3
		if pageContainsHeader {
			y1 = 27.5
			y2 = 22.5
		}

		// Draw Lines - left section
		pdf.SetDashPattern([]float64{}, 0)
		pdf.SetLineWidth(0.4)
		pdf.SetDrawColor(0, 0, 0)
		pdf.Line(col1.X, 19.3, col1.X, endY)
		pdf.SetLineWidth(0.3)
		pdf.SetDrawColor(160, 160, 160)
		pdf.SetDashPattern([]float64{0.5, 0.5}, 0)
		pdf.Line(col2.X, y1, col2.X, endY)
		pdf.Line(col3.X, y1, col3.X, endY)
		pdf.Line(col4.X, y1, col4.X, endY)

		// Draw Lines - middle section
		pdf.SetDashPattern([]float64{}, 0)
		pdf.SetLineWidth(0.4)
		pdf.Line(col5.X, y2, col5.X, endY)
		pdf.SetDashPattern([]float64{0.5, 0.5}, 0)
		pdf.SetLineWidth(0.3)
		pdf.Line(col6.X, y1, col6.X, endY)
		pdf.Line(col7.X, y1, col7.X, endY)
		pdf.Line(col8.X, y1, col8.X, endY)
		pdf.Line(col9.X, y1, col9.X, endY)
		pdf.SetDashPattern([]float64{}, 0)
		pdf.SetLineWidth(0.4)
		pdf.Line(col10.X, y2, col10.X, endY)

		// Draw Lines - right section
		pdf.SetDashPattern([]float64{0.5, 0.5}, 0)
		pdf.SetLineWidth(0.3)
		pdf.Line(col11.X, y1, col11.X, endY)
		pdf.Line(col12.X, y1, col12.X, endY)
		pdf.Line(col13.X, y1, col13.X, endY)
		pdf.Line(col14.X, y1, col14.X, endY)
		pdf.SetLineWidth(0.4)
		pdf.SetDrawColor(0, 0, 0)
		pdf.SetDashPattern([]float64{}, 0)
		pdf.Line(269.5, 19.3, 269.5, endY)
	}

	newPageIfNecessary := func(bottomY float64) {
		if bottomY > 196.7 {
			drawVerticalLines(y)
			pdf.AddPage()
			pageContainsHeader = false
			drawTopHorizontalLine()
			y = 19.3 // start Y at the top of the page
		}
	}

	fillHeader := func() {
		// Draw top line
		drawTopHorizontalLine()

		// Write the header line 1
		pdf.SetFillColor(155, 194, 230)
		pdf.SetFontSize(5)
		pdf.SetXY(col5.X, 19.7)
		pdf.CellFormat(91.4, 2.5, "Client new request", "", 0, "MC", true, 0, "")

		// Draw background for header line 2
		pdf.SetFillColor(191, 191, 191)
		pdf.Rect(col1.X, 22.5, 263, 5, "F")

		// Write the header line 2
		pdf.SetXY(col1.X, 22.5)
		pdf.CellFormat(col1.Wd, 5, "Block", "", 0, "MC", true, 0, "")
		pdf.SetXY(col2.X, 22.5)
		pdf.CellFormat(col2.Wd, 5, "Facade", "", 0, "MC", true, 0, "")
		pdf.SetXY(col3.X, 22.5)
		pdf.CellFormat(col3.Wd, 5, "Unit", "", 0, "MC", true, 0, "")
		pdf.SetTextColor(255, 0, 0)
		pdf.SetFontStyle("S")
		pdf.SetXY(col4.X, 22.5)
		pdf.CellFormat(col4.Wd, 5, "Original Request (Superseded)", "", 0, "MC", true, 0, "")
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFontStyle("")
		pdf.SetXY(col5.X, 22.5)
		pdf.CellFormat(col5.Wd, 5, "Aperture Reference", "", 0, "MC", true, 0, "")
		pdf.SetXY(col6.X, 22.5)
		pdf.CellFormat(col6.Wd, 5, "Type / System", "", 0, "MC", true, 0, "")
		pdf.SetXY(col7.X, 22.5)
		pdf.CellFormat(col7.Wd, 5, "Width (mm)", "", 0, "MC", true, 0, "")
		pdf.SetXY(col8.X, 22.5)
		pdf.CellFormat(col8.Wd, 5, "Height (mm)", "", 0, "MC", true, 0, "")
		pdf.SetXY(col9.X, 22.5)
		pdf.CellFormat(col9.Wd, 5, "Glass", "", 0, "MC", true, 0, "")
		pdf.SetXY(col10.X, 22.5)
		pdf.CellFormat(col10.Wd, 5, "UNIT", "", 0, "MC", true, 0, "")
		pdf.SetXY(col11.X, 22.5)
		pdf.CellFormat(col11.Wd, 5, "AMOUNT", "", 0, "MC", true, 0, "")
		pdf.SetXY(col12.X, 22.5)
		pdf.CellFormat(col12.Wd, 5, "", "", 0, "MC", true, 0, "")
		pdf.SetXY(col13.X, 22.5)
		pdf.CellFormat(col13.Wd, 5, "TOTAL AREA", "", 0, "MC", true, 0, "")
		pdf.SetXY(col14.X, 22.5)
		pdf.CellFormat(col14.Wd, 5, "", "", 0, "MC", true, 0, "")

		// Draw horizontal dashed-line at the bottom of header line 2
		pdf.SetDashPattern([]float64{0.5, 0.5}, 0)
		pdf.SetLineWidth(0.4)
		pdf.SetDrawColor(160, 160, 160)
		pdf.Line(col1.X, 27.5, 269.5, 27.5)
	}

	fillBody := func() {
		displayAperture := func(aperture *apertureEntity.ApertureEntity) {
			type Cell struct {
				Text        string
				Lines       [][]byte
				Wd          float64
				Ht          float64
				Align       string
				FontStyle   string
				R           int
				G           int
				B           int
				HasCurrency bool
			}

			cellList := [14]Cell{}
			// Fill in the cellList array
			{
				typeSystem := aperture.Profile.FriendlyName

				reference := fmt.Sprintf("%s-%d-%d-%s", typeSystem, aperture.ActualWidth, aperture.ActualHeight, aperture.Glass.GlazingStructure)

				cellList[0].Text = fmt.Sprintf("BLOCK %s", aperture.Block)
				cellList[0].Align = "C"

				cellList[1].Text = fmt.Sprintf("Facade %s", aperture.Facade)
				cellList[1].Align = "C"

				cellList[2].Text = fmt.Sprintf("%s", aperture.GetUnitFullName())
				cellList[2].Align = "C"

				cellList[3].Text = aperture.SupplierComment
				cellList[3].Align = "C"
				cellList[3].FontStyle = "S"
				cellList[3].R = 255

				cellList[4].Text = reference
				cellList[4].Align = "C"

				cellList[5].Text = typeSystem
				cellList[5].Align = "C"

				cellList[6].Text = formatters.ToReadableNumber_2dp(float64(aperture.ActualWidth))
				cellList[6].Align = "R"

				cellList[7].Text = formatters.ToReadableNumber_2dp(float64(aperture.ActualHeight))
				cellList[7].Align = "R"

				cellList[8].Text = aperture.Glass.GlazingStructure
				cellList[8].Align = "C"

				cellList[9].Text = "piece"
				cellList[9].Align = "C"

				cellList[10].Text = formatters.ToReadableWholeNumber(aperture.Quantity)
				cellList[10].Align = "C"

				totalArea := float64(aperture.Quantity) * float64(aperture.ActualWidth) * float64(aperture.ActualHeight) / 1000000.0
				cellList[12].Text = formatters.ToReadableNumber_2dp(totalArea)
				cellList[12].Align = "R"

				for i := range cellList {
					cellList[i].Wd = colArr[i].Wd
				}
			}

			y = pdf.GetY() // The current y within the pdf
			lineHt := 2.8

			maxHt := lineHt
			cellGapWd := 0.0
			cellGapHt := 1.0
			// Calculate Cell height
			for i := range 14 {
				cellList[i].Lines = pdf.SplitLines([]byte(cellList[i].Text), cellList[i].Wd-2*cellGapWd)
				cellList[i].Ht = float64(len(cellList[i].Lines)) * lineHt
				if cellList[i].Ht > maxHt {
					maxHt = cellList[i].Ht
				}
			}

			// Add new page if necessary
			{
				bottomY := y + maxHt + 2*cellGapHt
				newPageIfNecessary(bottomY)
			}

			// Render Cell
			x := col1.X // The current x within the pdf, starting as the left margin
			for i := range 14 {
				cell := cellList[i]
				pdf.SetTextColor(cell.R, cell.G, cell.B)
				pdf.SetFontStyle(cell.FontStyle)
				cellY := y + cellGapHt + (maxHt-cell.Ht)/2
				for splitJ := range cell.Lines {
					pdf.SetXY(x+cellGapWd, cellY)
					pdf.CellFormat(cell.Wd-2*cellGapWd, lineHt, string(cell.Lines[splitJ]), "", 0,
						cell.Align, false, 0, "")

					cellY += lineHt
				}
				x += cell.Wd
			}
			y += maxHt + 2*cellGapHt

			// Draw the bottom line for the current row
			pdf.SetLineWidth(0.3)
			pdf.SetDrawColor(160, 160, 160)
			pdf.SetDashPattern([]float64{0.5, 0.5}, 0)
			pdf.Line(col1.X, y, 269.5, y)

			pdf.SetY(y)
		}

		pdf.SetY(27.5) // start Y as the bottom of the header section
		for _, aperture := range mainApertureList {
			displayAperture(aperture)
		}
		drawVerticalLines(y)
	}

	fillFooter := func() {
		bottomY := y + 2.5
		newPageIfNecessary(bottomY)

		// Draw the horizontal line
		pdf.SetLineWidth(0.4)
		pdf.SetDrawColor(0, 0, 0)
		pdf.SetDashPattern([]float64{}, 0)
		pdf.Line(col1.X, y+2.5, 269.5, y+2.5)

		// Draw the vertical lines
		pdf.Line(col1.X, y, col1.X, y+2.5)
		pdf.Line(269.5, y, 269.5, y+2.5)
		pdf.SetDrawColor(160, 160, 160)
		pdf.Line(col5.X, y, col5.X, y+2.5)
		pdf.Line(col10.X, y, col10.X, y+2.5)
	}

	fillHeader()
	fillBody()
	fillFooter()

	var buff bytes.Buffer
	err := pdf.Output(&buff)
	if err != nil {
		return []byte{}, err
	}

	return buff.Bytes(), nil
}

func generateSection3(quote *quoteEntity.QuoteEntity, mainApertureList []*apertureEntity.ApertureEntity) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(5, 5, 5)
	pdf.AddPage()
	pdf.SetFont("Helvetica", "", 12)

	var y float64
	var aperturesInPage int // To keeps track of the amount of apertures in the current page

	fillPageHeader := func() {
		const (
			margin = 6.3
		)

		fillLeftSection := func() {
			const (
				x        float64 = margin
				imgY     float64 = 6.2
				imgWd    float64 = 45
				imgHt    float64 = 24
				lnHt     float64 = 1.4
				line1Y   float64 = imgY + imgHt + 1.2
				line2Y   float64 = line1Y + lnHt + 2.5
				line3Y   float64 = line2Y + lnHt + 3.1
				line4Y   float64 = line3Y + lnHt + 3.1
				vatLineY float64 = line4Y + lnHt + 3.1
			)

			// Draw the logo
			var opt fpdf.ImageOptions
			opt.ImageType = "png"
			pdf.ImageOptions(pathUtils.ToImagePath("luxe-logo.PNG"), x, imgY, imgWd, imgHt, false, opt, 0, "")

			// Line 1
			pdf.SetFontSize(8.2)
			pdf.SetFontStyle("B")
			pdf.SetTextColor(26, 62, 99)
			pdf.SetXY(x, line1Y)
			pdf.CellFormat(100, 2, "LUXE Apertures Ltd.", "", 0, "TL", false, 0, "")

			// Line 2
			pdf.SetFontSize(9)
			pdf.SetFontStyle("")
			pdf.SetTextColor(140, 140, 140)
			pdf.SetXY(x, line2Y)
			pdf.CellFormat(100, 2, "5, 'Villa Romana'", "", 0, "TL", false, 0, "")

			// Line 3
			pdf.SetXY(x, line3Y)
			pdf.CellFormat(100, 2, "Xrobb l-Ghagin, Limiti ta'", "", 0, "TL", false, 0, "")

			// Line 4
			pdf.SetXY(x, line4Y)
			pdf.CellFormat(100, 2, "Marsaxlokk, MXK4080", "", 0, "TL", false, 0, "")

			// Line 4
			pdf.SetXY(x, vatLineY)
			pdf.CellFormat(100, 2, "Vat: MT2902-3304", "", 0, "TL", false, 0, "")
		}

		fillRightSection := func() {
			const (
				x      = 135.5
				lineHt = 4
			)

			pdf.SetFontSize(7.9)

			// Client name Line
			y = 26.5
			clientFullname := quote.Client.GetFullName()
			pdf.SetXY(x, y)
			pdf.CellFormat(100, 2, "Customer : "+clientFullname, "", 0, "TL", false, 0, "")

			// Address Line
			addressPieces := strings.Split(quote.Client.Address, "\n")
			address := ""
			if len(addressPieces) > 0 {
				address = addressPieces[0]
			}
			if len(addressPieces) > 1 {
				address += " - " + addressPieces[1]
			}
			for i := 2; i < len(addressPieces); i++ {
				address += ", " + addressPieces[i]
			}

			projectWidth := pdf.GetStringWidth("Project : ")
			addressLines := pdf.SplitLines([]byte(address), 210-margin-x-projectWidth)

			y += lineHt
			pdf.SetXY(x, y)
			pdf.CellFormat(projectWidth, 2, "Project : ", "", 0, "TL", false, 0, "")

			if len(addressLines) > 0 {
				pdf.SetXY(x+projectWidth, y)
				pdf.CellFormat(100, 2, string(addressLines[0]), "", 0, "TL", false, 0, "")

				for i := 1; i < len(addressLines); i++ {
					y += lineHt
					pdf.SetXY(x+projectWidth, y)
					pdf.CellFormat(100, 2, string(addressLines[i]), "", 0, "TL", false, 0, "")
				}
			}

			// Date Line
			y += lineHt
			now := time.Now()
			pdf.SetXY(x, y)
			pdf.CellFormat(5, 2, "Date : "+now.Format("02.01.2006"), "", 0, "TL", false, 0, "")

			// Quote line
			y += lineHt
			pdf.SetXY(x, y)
			pdf.CellFormat(5, 2, fmt.Sprintf("Quote : %d", quote.GetNumber()), "", 0, "TL", false, 0, "")
		}

		fillLeftSection()
		fillRightSection()
	}

	newPage := func() {
		pdf.AddPage()
		fillPageHeader()
		y = 55.5
		aperturesInPage = 0
	}

	fillAperture := func(aperture *apertureEntity.ApertureEntity, childApertureNo int) {
		fillTopSection := func() {
			const (
				margin      = 4
				hLineGap    = 4.2 // The gap between the two horizontal lines
				middleLineX = 105.2
			)

			if len(aperture.ChildApertures) > 0 {
				// Draw background for the header line
				pdf.SetFillColor(155, 194, 230)
				pdf.Rect(margin, y, 210-margin*2, 5, "F")

				// Write the header line
				text := "Combined Aperture"
				textWd := pdf.GetStringWidth(text)
				pdf.SetTextColor(0, 0, 0)
				pdf.SetFontStyle("B")
				pdf.SetFontSize(11)
				pdf.SetXY((210-margin*2-textWd)/2, y)
				pdf.CellFormat(20, 5, text, "", 0, "ML", true, 0, "")
				y += 5
			} else if childApertureNo > 0 {
				// Draw background for the header line
				pdf.SetFillColor(155, 194, 230)
				pdf.Rect(margin, y, 210-margin*2, 5, "F")

				// Write the header line
				text := fmt.Sprintf("Component Aperture %d", childApertureNo)
				textWd := pdf.GetStringWidth(text)
				pdf.SetTextColor(0, 0, 0)
				pdf.SetFontStyle("B")
				pdf.SetFontSize(11)
				pdf.SetXY((210-margin*2-textWd)/2, y)
				pdf.CellFormat(20, 5, text, "", 0, "ML", true, 0, "")
				y += 5
			}

			pdf.SetLineWidth(0.3)
			pdf.SetDrawColor(77, 77, 79)

			// Draw Top Horizontal Lines
			pdf.Line(margin, y, 210-margin, y)
			pdf.Line(margin, y+hLineGap, 210-margin, y+hLineGap)

			// Draw Top Vertical Lines
			pdf.Line(margin, y, margin, y+hLineGap)
			pdf.Line(middleLineX, y, middleLineX, y+hLineGap)
			pdf.Line(210-margin, y, 210-margin, y+hLineGap)

			// Write Titles
			pdf.SetFontSize(9.5)
			pdf.SetTextColor(88, 144, 82)
			pdf.SetFontStyle("B")
			pdf.SetXY(margin, y)
			pdf.CellFormat(middleLineX-margin, hLineGap, fmt.Sprintf("Unit : %s", aperture.GetUnitFullName()), "", 0, "CM", false, 0, "")

			pdf.SetXY(middleLineX, y)
			pdf.CellFormat(210-margin-middleLineX, hLineGap, fmt.Sprintf("Dimensions (W x H) : %dx%d mm", aperture.ActualWidth, aperture.ActualHeight), "", 0, "CM", false, 0, "")

			y += hLineGap
		}

		leftSectionBottomY := y
		rightSectionBottomY := y

		fillLeftSection := func() {
			const (
				imgTopGap = 13
				maxImgWd  = 73.0
				maxImgHt  = 28.0
			)

			// Write the text
			pdf.SetFontSize(6)
			pdf.SetTextColor(77, 77, 79)
			pdf.SetXY(21.6, y+3.5)
			pdf.CellFormat(5, 2, fmt.Sprintf("Unit:%s / Pcs:%d", aperture.GetUnitFullName(), aperture.GetQuantity()), "", 0, "TL", false, 0, "")

			// Draw the image together with dimensions
			ext := ""
			switch fileTypes.FileType(aperture.Image.Type) {
			case fileTypes.Png:
				ext = "png"
			case fileTypes.Jpg:
				ext = "jpg"
			case fileTypes.Gif:
				ext = "gif"
			}

			imgName := fmt.Sprintf("%d.%s", aperture.ID, ext)

			reader := bytes.NewReader(aperture.Image.Data)
			pdf.RegisterImageOptionsReader(imgName, fpdf.ImageOptions{ImageType: ext}, reader)

			imgWd := maxImgWd
			imgHt := maxImgHt
			info := pdf.GetImageInfo(imgName)
			if info != nil {
				imgWd = info.Width()
				imgHt = info.Height()

				if imgWd > maxImgWd && imgHt <= maxImgHt {
					ratio := maxImgWd / imgWd
					imgWd = maxImgWd
					imgHt *= ratio
				} else if imgWd <= maxImgWd && imgHt > maxImgHt {
					ratio := maxImgHt / imgHt
					imgHt = maxImgHt
					imgWd *= ratio
				} else if imgWd > maxImgWd && imgHt > maxImgHt {
					wdRatio := maxImgWd / imgWd
					htRatio := maxImgHt / imgHt
					ratio := math.Min(wdRatio, htRatio)
					imgWd *= ratio
					imgHt *= ratio
				}
			}
			pdf.ImageOptions(imgName, 21.6, y+imgTopGap, imgWd, imgHt, false, fpdf.ImageOptions{}, 0, "")

			leftSectionBottomY = y + imgTopGap + imgHt

			if !aperture.IsCombined {
				// Draw the left-hand side of the image
				{
					const (
						x = 19.6
					)

					startY := y + imgTopGap
					endY := y + imgTopGap + imgHt

					pdf.SetLineWidth(0.4)

					// Draw the small horizontal lines
					pdf.Line(x-2, startY, x+2, startY)
					pdf.Line(x-2, endY, x+2, endY)

					// Draw the vertical line
					pdf.SetLineWidth(0.3)
					pdf.Line(x, startY, x, endY)

					// Draw the red triangles
					pdf.SetFillColor(237, 77, 58)
					pdf.Polygon([]fpdf.PointType{{X: x, Y: startY}, {X: x + 1.5, Y: startY + 3}, {X: x, Y: startY + 2}, {X: x - 1.5, Y: startY + 3}}, "F")
					pdf.Polygon([]fpdf.PointType{{X: x, Y: endY}, {X: x + 1.5, Y: endY - 3}, {X: x, Y: endY - 2}, {X: x - 1.5, Y: endY - 3}}, "F")

					// Write the dimension
					txt := strconv.Itoa(aperture.ActualHeight)
					pdf.SetFontSize(9)

					txtWd := pdf.GetStringWidth(txt)
					_, txtHt := pdf.GetFontSize()
					txtX := x - txtWd/2
					txtY := startY + (endY-startY-txtHt)/2

					pdf.TransformBegin()
					pdf.TransformRotate(90, txtX+txtWd/2, txtY+txtHt/2)
					pdf.SetFillColor(255, 255, 255)
					pdf.SetFontStyle("B")
					pdf.SetXY(txtX, txtY)
					pdf.CellFormat(txtWd, txtHt, txt, "", 0, "CM", true, 0, "")
					pdf.TransformEnd()
				}

				// Draw the bottom side of the image
				{
					y := leftSectionBottomY + 3
					startX := 21.6
					endX := 21.6 + imgWd

					pdf.SetLineWidth(0.4)

					// Draw the small vertical lines
					pdf.Line(startX, y-2, startX, y+2)
					pdf.Line(endX, y-2, endX, y+2)

					// Draw the horizontal line
					pdf.SetLineWidth(0.3)
					pdf.Line(startX, y, endX, y)

					// Draw the red triangles
					pdf.SetFillColor(237, 77, 58)
					pdf.Polygon([]fpdf.PointType{{X: startX, Y: y}, {X: startX + 3, Y: y + 1.5}, {X: startX + 2, Y: y}, {X: startX + 3, Y: y - 1.5}}, "F")
					pdf.Polygon([]fpdf.PointType{{X: endX, Y: y}, {X: endX - 3, Y: y + 1.5}, {X: endX - 2, Y: y}, {X: endX - 3, Y: y - 1.5}}, "F")

					// Write the dimension
					txt := strconv.Itoa(aperture.ActualWidth)
					pdf.SetFontSize(7)

					txtWd := pdf.GetStringWidth(txt)
					_, txtHt := pdf.GetFontSize()
					txtX := startX + (endX-startX-txtWd)/2
					txtY := y - txtHt/2

					pdf.SetFillColor(255, 255, 255)
					pdf.SetFontStyle("B")
					pdf.SetXY(txtX, txtY)
					pdf.CellFormat(txtWd, txtHt, txt, "", 0, "CM", true, 0, "")

					leftSectionBottomY = y + 2
				}
			}
		}

		fillRightSection := func() {
			const (
				col1X     = 100
				col1Width = 31
				col2X     = col1X + col1Width + 1
				col2Width = 70
				lineHt    = 3.5
			)

			// Brand Line
			pdf.SetFontSize(7.7)
			pdf.SetFontStyle("")
			y += 3.5
			pdf.SetXY(col1X, y)
			pdf.CellFormat(col1Width, 2, "Brand :", "", 0, "TR", false, 0, "")
			pdf.SetXY(col2X, y)
			pdf.CellFormat(col2Width, 2, aperture.Profile.Brand, "", 0, "TL", false, 0, "")

			// Serie Line
			y += lineHt
			pdf.SetXY(col1X, y)
			pdf.CellFormat(col1Width, 2, "Series :", "", 0, "TR", false, 0, "")
			pdf.SetXY(col2X, y)
			pdf.CellFormat(col2Width, 2, aperture.Profile.Serie, "", 0, "TL", false, 0, "")

			// System Type
			y += lineHt
			pdf.SetXY(col1X, y)
			pdf.CellFormat(col1Width, 2, "System Type :", "", 0, "TR", false, 0, "")
			pdf.SetXY(col2X, y)
			pdf.CellFormat(col2Width, 2, aperture.FriendlyName, "", 0, "TL", false, 0, "")

			// Profile Thickness
			y += lineHt
			pdf.SetXY(col1X, y)
			pdf.CellFormat(col1Width, 2, "Profile Thickness :", "", 0, "TR", false, 0, "")
			pdf.SetXY(col2X, y)
			pdf.CellFormat(col2Width, 2, fmt.Sprintf("%.0f mm", aperture.Profile.SystemThickness), "", 0, "TL", false, 0, "")

			// Material
			y += lineHt
			pdf.SetXY(col1X, y)
			pdf.CellFormat(col1Width, 2, "Material :", "", 0, "TR", false, 0, "")
			pdf.SetXY(col2X, y)
			pdf.CellFormat(col2Width, 2, aperture.Profile.Material, "", 0, "TL", false, 0, "")

			// Color
			y += lineHt
			pdf.SetXY(col1X, y)
			pdf.CellFormat(col1Width, 2, "Color :", "", 0, "TR", false, 0, "")
			pdf.SetXY(col2X, y)
			pdf.CellFormat(col2Width, 2, aperture.Color.FriendlyName, "", 0, "TL", false, 0, "")

			// Glass
			y += lineHt
			pdf.SetXY(col1X, y)
			pdf.CellFormat(col1Width, 2, "Glass :", "", 0, "TR", false, 0, "")
			pdf.SetXY(col2X, y)
			pdf.CellFormat(col2Width, 2, aperture.Glass.GlazingStructure, "", 0, "TL", false, 0, "")

			// Accessories
			if aperture.Accessories != "" {
				accessoryPieces := strings.Split(aperture.Accessories, "\n")

				for _, accessory := range accessoryPieces {
					pdf.SetXY(col1X, y+lineHt)
					pdf.CellFormat(col1Width, 2, "Accessory :", "", 0, "TR", false, 0, "")

					accessoryLines := pdf.SplitLines([]byte(accessory), col2Width)
					for _, accessoryLine := range accessoryLines {
						y += lineHt
						pdf.SetXY(col2X, y)
						pdf.CellFormat(col2Width, 2, string(accessoryLine), "", 0, "TL", false, 0, "")
					}
				}
			}

			// Area
			y += lineHt
			area := float64(aperture.ActualWidth) * float64(aperture.ActualHeight) / 1000000.0
			pdf.SetXY(col1X, y)
			pdf.CellFormat(col1Width, 2, "Total Area (SQM) :", "", 0, "TR", false, 0, "")
			pdf.SetXY(col2X, y)
			pdf.CellFormat(col2Width, 2, fmt.Sprintf("%.2f", area), "", 0, "TL", false, 0, "")

			// Glass U-Value
			y += lineHt
			pdf.SetXY(col1X, y)
			pdf.CellFormat(col1Width, 2, "Insulation Val-Glass :", "", 0, "TR", false, 0, "")
			text := strings.Replace(fmt.Sprintf("%.2f W/m", aperture.Glass.UValue), ".", ",", 1)
			textWd := pdf.GetStringWidth(text)
			pdf.SetXY(col2X, y)
			pdf.CellFormat(textWd, 2, text, "", 0, "TL", false, 0, "")
			pdf.SetFontSize(5)
			textWd = pdf.GetStringWidth("2")
			pdf.CellFormat(textWd, 2, "2", "", 0, "TL", false, 0, "")
			pdf.SetFontSize(7.7)
			pdf.CellFormat(textWd, 2, ".K", "", 0, "TL", false, 0, "")

			// Profile U-Value
			y += lineHt
			pdf.SetXY(col1X, y)
			pdf.CellFormat(col1Width, 2, "Insulation Val-Profile :", "", 0, "TR", false, 0, "")
			text = strings.Replace(fmt.Sprintf("%.2f W/m", aperture.Profile.UValue), ".", ",", 1)
			textWd = pdf.GetStringWidth(text)
			pdf.SetXY(col2X, y)
			pdf.CellFormat(textWd, 2, text, "", 0, "TL", false, 0, "")
			pdf.SetFontSize(5)
			textWd = pdf.GetStringWidth("2")
			pdf.CellFormat(textWd, 2, "2", "", 0, "TL", false, 0, "")
			pdf.SetFontSize(7.7)
			pdf.CellFormat(textWd, 2, ".K", "", 0, "TL", false, 0, "")

			// Overall U-Value
			y += lineHt
			pdf.SetXY(col1X, y)
			pdf.CellFormat(col1Width, 2, "Insulation Val-Product :", "", 0, "TR", false, 0, "")
			text = strings.Replace(fmt.Sprintf("%.2f W/m", aperture.OverallUValue), ".", ",", 1)
			textWd = pdf.GetStringWidth(text)
			pdf.SetFontStyle("B")
			pdf.SetXY(col2X, y)
			pdf.CellFormat(textWd, 2, text, "", 0, "TL", false, 0, "")
			pdf.SetFontSize(5)
			textWd = pdf.GetStringWidth("2")
			pdf.CellFormat(textWd, 2, "2", "", 0, "TL", false, 0, "")
			pdf.SetFontSize(7.7)
			pdf.CellFormat(textWd, 2, ".K", "", 0, "TL", false, 0, "")
			pdf.SetFontStyle("")

			// Note
			if aperture.SupplierComment != "" {
				pdf.SetXY(col1X, y+lineHt)
				pdf.CellFormat(col1Width, 2, "Note :", "", 0, "TR", false, 0, "")

				noteLines := pdf.SplitLines([]byte(aperture.SupplierComment), col2Width)
				for _, noteLine := range noteLines {
					y += lineHt
					pdf.SetXY(col2X, y)
					pdf.CellFormat(col2Width, 2, string(noteLine), "", 0, "TL", false, 0, "")
				}
			}

			rightSectionBottomY = y + lineHt
		}

		fillBottomSection := func() {
			const (
				margin   = 4
				hLineGap = 4.2 // The gap between the two horizontal lines
			)

			y += 17

			// Draw Horizontal Lines
			pdf.Line(margin, y, 210-margin, y)
			pdf.Line(margin, y+hLineGap, 210-margin, y+hLineGap)

			// Draw Vertical Lines
			pdf.Line(margin, y, margin, y+hLineGap)
			pdf.Line(margin+67, y, margin+67, y+hLineGap)
			pdf.Line(margin+135, y, margin+135, y+hLineGap)
			pdf.Line(210-margin, y, 210-margin, y+hLineGap)

			// Quantity Line
			pdf.SetTextColor(72, 68, 135)
			pdf.SetFontStyle("B")
			pdf.SetXY(margin+67, y)
			pdf.CellFormat(67, hLineGap, fmt.Sprintf("Quantity : %d,00", aperture.Quantity), "", 0, "CM", false, 0, "")

			// Draw the Bottom Line
			y += hLineGap + 2.8
			pdf.SetLineWidth(1)
			pdf.Line(margin, y, 210-margin, y)
		}

		fillTopSection()
		fillLeftSection()
		fillRightSection()

		y = leftSectionBottomY
		if rightSectionBottomY > leftSectionBottomY {
			y = rightSectionBottomY
		}

		fillBottomSection()
	}

	fillPageHeader()
	y = 55.5

	for _, mainAperture := range mainApertureList {
		if aperturesInPage < 2 {
			y += 2.2
		} else {
			newPage()
		}

		fillAperture(mainAperture, 0)
		aperturesInPage++

		for c, childAperture := range mainAperture.ChildApertures {
			childAperture.SetUnitName(mainAperture.UnitName)
			childAperture.SetUnitNumber(mainAperture.UnitNumber)

			if aperturesInPage < 2 {
				y += 2.2
			} else {
				newPage()
			}

			fillAperture(childAperture, c+1)
			aperturesInPage++
		}
	}

	var buff bytes.Buffer
	err := pdf.Output(&buff)
	if err != nil {
		return []byte{}, err
	}

	return buff.Bytes(), nil
}
