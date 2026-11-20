package supplierPdf

import (
	_ "embed"

	"codeberg.org/go-pdf/fpdf"
	"github.com/pkg/errors"
)

// Language is the language the fixed texts of the supplier PDF are written in.
type Language string

const (
	// LanguageEnglish is the language of the report as it has always been.
	LanguageEnglish Language = "en"
	// LanguageTurkish translates the fixed texts of the report to Turkish.
	LanguageTurkish Language = "tr"
)

// ParseLanguage converts the text "en" or "tr" into a Language.
func ParseLanguage(text string) (Language, error) {
	switch Language(text) {
	case LanguageEnglish:
		return LanguageEnglish, nil
	case LanguageTurkish:
		return LanguageTurkish, nil
	}

	return "", errors.Errorf("'%s' is not a valid language. It must be either en or tr", text)
}

// PageNumberText is the text stamped on every page of the final document. %p is the page number, %P the number of pages.
func (l Language) PageNumberText() string {
	return textsFor(l).PageNumber
}

// The core PDF fonts cannot draw some Turkish letters (ğ, ş, ı, İ), so Turkish reports use an embedded font.
// Liberation Sans has the same letter widths as Helvetica/Arial, so the layout stays the same.
// Its licence (SIL Open Font License 1.1) is in fonts/LICENSE-LiberationSans.txt.
const unicodeFontFamily = "LiberationSans"

//go:embed fonts/LiberationSans-Regular.ttf
var liberationSansRegular []byte

//go:embed fonts/LiberationSans-Bold.ttf
var liberationSansBold []byte

// prepareFonts registers the fonts needed for the language and returns the name of the font family to use.
func prepareFonts(pdf *fpdf.Fpdf, language Language) string {
	if language != LanguageTurkish {
		return "Helvetica"
	}

	pdf.AddUTF8FontFromBytes(unicodeFontFamily, "", liberationSansRegular)
	pdf.AddUTF8FontFromBytes(unicodeFontFamily, "B", liberationSansBold)
	return unicodeFontFamily
}

// texts holds every fixed text of sections 2 and 3 of the supplier PDF. Texts that contain %s or %d are format strings.
// Texts that come from the database (names, descriptions, comments) are not translated.
type texts struct {
	// Section 2
	ClientNewRequest          string
	Block                     string
	BlockValue                string // format: block
	Room                      string
	Unit                      string
	OriginalRequestSuperseded string
	ApertureReference         string
	TypeSystem                string
	WidthMm                   string
	HeightMm                  string
	Glass                     string
	UnitUpper                 string
	AmountUpper               string
	TotalAreaUpper            string
	Piece                     string

	// Section 3: page header
	VatLine  string
	Customer string
	Project  string
	Date     string
	Quote    string // format: quote number

	// Section 3: apertures
	CombinedAperture  string
	ComponentAperture string // format: number of the component
	Sku               string // format: unit name, id
	Dimensions        string // format: width, height
	UnitPieces        string // format: unit name, quantity
	InteriorView      string
	Brand             string
	Series            string
	SystemType        string
	ProfileThickness  string
	Material          string
	Color             string
	GlassLabel        string
	Accessory         string
	InsectScreen      string
	NotIncluded       string
	Included          string
	TotalAreaSqm      string
	InsulationGlass   string
	InsulationProfile string
	InsulationProduct string
	Note              string
	Quantity          string // format: quantity
	Summary           string
	AmountOfApertures string
	PageNumber        string
}

var englishTexts = texts{
	ClientNewRequest:          "Client new request",
	Block:                     "Block",
	BlockValue:                "BLOCK %s",
	Room:                      "Room",
	Unit:                      "Unit",
	OriginalRequestSuperseded: "Original Request (Superseded)",
	ApertureReference:         "Aperture Reference",
	TypeSystem:                "Type / System",
	WidthMm:                   "Width (mm)",
	HeightMm:                  "Height (mm)",
	Glass:                     "Glass",
	UnitUpper:                 "UNIT",
	AmountUpper:               "AMOUNT",
	TotalAreaUpper:            "TOTAL AREA",
	Piece:                     "piece",

	VatLine:  "Vat: MT2902-3304",
	Customer: "Customer : ",
	Project:  "Project : ",
	Date:     "Date : ",
	Quote:    "Quote : %d",

	CombinedAperture:  "Combined Aperture",
	ComponentAperture: "Component Aperture %d",
	Sku:               "SKU : %s - %d",
	Dimensions:        "Dimensions (W x H) : %dx%d mm",
	UnitPieces:        "Unit:%s / Pcs:%d",
	InteriorView:      "Interior View",
	Brand:             "Brand :",
	Series:            "Series :",
	SystemType:        "System Type :",
	ProfileThickness:  "Profile Thickness :",
	Material:          "Material :",
	Color:             "Color :",
	GlassLabel:        "Glass :",
	Accessory:         "Accessory :",
	InsectScreen:      "Insect screen :",
	NotIncluded:       "Not included",
	Included:          "Included",
	TotalAreaSqm:      "Total Area (SQM) :",
	InsulationGlass:   "Insulation Val-Glass :",
	InsulationProfile: "Insulation Val-Profile :",
	InsulationProduct: "Insulation Val-Product :",
	Note:              "Note :",
	Quantity:          "Quantity : %d,00",
	Summary:           "Summary",
	AmountOfApertures: "Amount of Apertures :",
	PageNumber:        "Page %p of %P",
}

var turkishTexts = texts{
	ClientNewRequest:          "Müşteri yeni talebi",
	Block:                     "Blok",
	BlockValue:                "BLOK %s",
	Room:                      "Oda",
	Unit:                      "Birim",
	OriginalRequestSuperseded: "Orijinal Talep (Değiştirildi)",
	ApertureReference:         "Doğrama Referansı",
	TypeSystem:                "Tip / Sistem",
	WidthMm:                   "Genişlik (mm)",
	HeightMm:                  "Yükseklik (mm)",
	Glass:                     "Cam",
	UnitUpper:                 "BİRİM",
	AmountUpper:               "ADET",
	TotalAreaUpper:            "TOPLAM ALAN",
	Piece:                     "adet",

	VatLine:  "KDV No: MT2902-3304",
	Customer: "Müşteri : ",
	Project:  "Proje : ",
	Date:     "Tarih : ",
	Quote:    "Teklif : %d",

	CombinedAperture:  "Birleşik Doğrama",
	ComponentAperture: "Bileşen Doğrama %d",
	Sku:               "SKU : %s - %d",
	Dimensions:        "Ölçüler (G x Y) : %dx%d mm",
	UnitPieces:        "Birim:%s / Adet:%d",
	InteriorView:      "İç Görünüm",
	Brand:             "Marka :",
	Series:            "Seri :",
	SystemType:        "Sistem Tipi :",
	ProfileThickness:  "Profil Kalınlığı :",
	Material:          "Malzeme :",
	Color:             "Renk :",
	GlassLabel:        "Cam :",
	Accessory:         "Aksesuar :",
	InsectScreen:      "Sineklik :",
	NotIncluded:       "Dahil değil",
	Included:          "Dahil",
	TotalAreaSqm:      "Toplam Alan (m²) :",
	InsulationGlass:   "Yalıtım Değeri-Cam :",
	InsulationProfile: "Yalıtım Değeri-Profil :",
	InsulationProduct: "Yalıtım Değeri-Ürün :",
	Note:              "Not :",
	Quantity:          "Adet : %d,00",
	Summary:           "Özet",
	AmountOfApertures: "Doğrama Adedi :",
	PageNumber:        "Sayfa %p / %P",
}

// textsFor returns the fixed texts of the given language. Anything other than Turkish gives the English texts.
func textsFor(language Language) texts {
	if language == LanguageTurkish {
		return turkishTexts
	}

	return englishTexts
}
