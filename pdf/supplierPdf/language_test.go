package supplierPdf

import (
	"Quote/db/entity/entities/apertureEntity"
	"Quote/db/entity/entities/colorEntity"
	"Quote/db/entity/entities/glassEntity"
	"Quote/db/entity/entities/profileEntity"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestParseLanguage(t *testing.T) {
	for text, want := range map[string]Language{"en": LanguageEnglish, "tr": LanguageTurkish} {
		if got, err := ParseLanguage(text); err != nil || got != want {
			t.Fatalf("ParseLanguage(%q) = %q, %v", text, got, err)
		}
	}
	for _, text := range []string{"", "EN", "Tr", "de", "english"} {
		if _, err := ParseLanguage(text); err == nil {
			t.Fatalf("ParseLanguage(%q) should fail", text)
		}
	}
}

// Every text must be filled in for both languages, and a translation must use the same format verbs as the English text,
// otherwise a Sprintf with the same arguments would print garbage.
func TestTextsAreComplete(t *testing.T) {
	verbs := regexp.MustCompile(`%[a-zA-Z]`)
	english, turkish := reflect.ValueOf(englishTexts), reflect.ValueOf(turkishTexts)

	for i := 0; i < english.NumField(); i++ {
		name := english.Type().Field(i).Name
		en, tr := english.Field(i).String(), turkish.Field(i).String()
		if strings.TrimSpace(en) == "" {
			t.Errorf("English text %s is empty", name)
		}
		if strings.TrimSpace(tr) == "" {
			t.Errorf("Turkish text %s is empty", name)
		}
		if !reflect.DeepEqual(verbs.FindAllString(en, -1), verbs.FindAllString(tr, -1)) {
			t.Errorf("%s: format verbs differ between %q and %q", name, en, tr)
		}
	}
}

func TestPageNumberText(t *testing.T) {
	for _, language := range []Language{LanguageEnglish, LanguageTurkish} {
		text := language.PageNumberText()
		if !strings.Contains(text, "%p") || !strings.Contains(text, "%P") {
			t.Fatalf("%s page number text %q must contain %%p and %%P", language, text)
		}
		// pdfcpu stamps the text with a standard PDF font, which cannot draw letters like ğ or ş
		for _, r := range text {
			if r > 0x7e {
				t.Fatalf("%s page number text %q has a character the standard font cannot draw: %q", language, text, r)
			}
		}
	}
	if LanguageEnglish.PageNumberText() != "Page %p of %P" {
		t.Fatalf("the English page numbers changed: %q", LanguageEnglish.PageNumberText())
	}
}

// Section 2 in Turkish: the letters that the standard PDF font cannot draw must come out right.
func TestSection2InTurkish(t *testing.T) {
	if _, err := exec.LookPath("pdftotext"); err != nil {
		t.Skip("pdftotext is not available")
	}

	aperture := apertureEntity.New()
	aperture.ID, aperture.Quantity, aperture.ActualWidth, aperture.ActualHeight = 1, 1, 1000, 2000
	aperture.Profile = &profileEntity.ProfileEntity{FriendlyName: "Profile"}
	aperture.Glass = &glassEntity.GlassEntity{GlazingStructure: "4-12-4"}
	aperture.Color = &colorEntity.ColorEntity{FriendlyName: "White"}
	list := []*apertureEntity.ApertureEntity{aperture}

	text := func(language Language) string {
		out, err := generateSection2(list, language)
		if err != nil {
			t.Fatal(err)
		}
		file := filepath.Join(t.TempDir(), "section2.pdf")
		if err := os.WriteFile(file, out, 0644); err != nil {
			t.Fatal(err)
		}
		extracted, err := exec.Command("pdftotext", "-layout", file, "-").Output()
		if err != nil {
			t.Fatal(err)
		}
		return string(extracted)
	}

	turkish := text(LanguageTurkish)
	for _, want := range []string{"Müşteri yeni talebi", "Doğrama Referansı", "Genişlik (mm)", "Yükseklik (mm)", "BİRİM", "TOPLAM ALAN", "adet"} {
		if !strings.Contains(turkish, want) {
			t.Errorf("the Turkish section 2 does not contain %q:\n%s", want, turkish)
		}
	}
	if strings.Contains(turkish, "Width (mm)") || strings.Contains(turkish, "Aperture Reference") {
		t.Errorf("the Turkish section 2 still has English headings:\n%s", turkish)
	}

	english := text(LanguageEnglish)
	for _, want := range []string{"Client new request", "Aperture Reference", "Width (mm)", "TOTAL AREA", "piece"} {
		if !strings.Contains(english, want) {
			t.Errorf("the English section 2 does not contain %q:\n%s", want, english)
		}
	}
}
