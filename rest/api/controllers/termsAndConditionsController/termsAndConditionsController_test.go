package termsAndConditionsController

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"

	"github.com/gin-gonic/gin"
)

func upload(t *testing.T, contentType string, content string, withAuth bool) int {
	t.Helper()
	gin.SetMode(gin.TestMode)

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	h := textproto.MIMEHeader{}
	h.Set("Content-Disposition", `form-data; name="file"; filename="t.pdf"`)
	h.Set("Content-Type", contentType)
	part, _ := w.CreatePart(h)
	part.Write([]byte(content))
	w.Close()

	r := gin.New()
	r.POST("/upload", Upload)
	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	if withAuth {
		req.SetBasicAuth("1", "x")
	}
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec.Code
}

func TestUploadRejectsInvalidFiles(t *testing.T) {
	if code := upload(t, "application/pdf", "%PDF-1.4", false); code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without auth, got %d", code)
	}
	if code := upload(t, "image/png", "%PDF-1.4", true); code != http.StatusUnsupportedMediaType {
		t.Fatalf("expected 415 for wrong content type, got %d", code)
	}
	if code := upload(t, "application/pdf", "not a pdf", true); code != http.StatusUnsupportedMediaType {
		t.Fatalf("expected 415 for wrong signature, got %d", code)
	}
}
