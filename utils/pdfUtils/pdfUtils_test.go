package pdfUtils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"codeberg.org/go-pdf/fpdf"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func makePdf(orientation, size string, pages int) []byte {
	pdf := fpdf.New(orientation, "mm", size, "")
	pdf.SetFont("Arial", "", 12)
	for i := 0; i < pages; i++ {
		pdf.AddPage()
		pdf.Cell(40, 10, "content")
	}
	var b bytes.Buffer
	pdf.Output(&b)
	return b.Bytes()
}

func TestAddPageNumbers(t *testing.T) {
	// Mixed orientations / sizes, as in the supplier report
	merged, err := MergePdfBytes(makePdf("L", "A4", 2), makePdf("P", "A4", 3), makePdf("P", "A5", 1))
	if err != nil {
		t.Fatal(err)
	}

	numbered, err := AddPageNumbers(merged)
	if err != nil {
		t.Fatal(err)
	}

	total, err := api.PageCount(bytes.NewReader(numbered), nil)
	if err != nil || total != 6 {
		t.Fatalf("expected 6 pages, got %d (%v)", total, err)
	}

	// Every page must carry the page number stamp
	dir := t.TempDir()
	if err := api.ExtractContent(bytes.NewReader(numbered), dir, "n", nil, nil); err != nil {
		t.Fatal(err)
	}
	for page := 1; page <= 6; page++ {
		data, err := os.ReadFile(filepath.Join(dir, fmt.Sprintf("n_Content_page_%d.txt", page)))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(data), "/Type /Pagination") {
			t.Fatalf("page %d has no page number stamp", page)
		}
	}

	// When a PDF text extractor is available, also check the stamped text
	if _, err := exec.LookPath("pdftotext"); err == nil {
		file := filepath.Join(dir, "numbered.pdf")
		if err := os.WriteFile(file, numbered, 0644); err != nil {
			t.Fatal(err)
		}
		for page := 1; page <= 6; page++ {
			out, err := exec.Command("pdftotext", "-f", fmt.Sprint(page), "-l", fmt.Sprint(page), file, "-").Output()
			if err != nil {
				t.Fatal(err)
			}
			if want := fmt.Sprintf("Page %d of 6", page); !strings.Contains(string(out), want) {
				t.Fatalf("page %d text %q does not contain %q", page, out, want)
			}
		}
	}
}
