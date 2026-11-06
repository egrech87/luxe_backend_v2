package apertureEntity

import (
	"Quote/utils/timeUtils"
	"testing"
	"time"
)

func TestActualDimensionsRoundTrip(t *testing.T) {
	if _, err := timeUtils.LoadLocation("UTC"); err != nil {
		t.Fatal(err)
	}
	e := New()
	e.CreatedAt, e.UpdatedAt = time.Now(), time.Now()
	err := e.FillFromJson(map[string]any{
		"quantity": 2, "actual_width": []any{float64(900), float64(905)}, "actual_height": "{1200,1200}",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(e.ActualWidth) != 2 || e.ActualHeight[1] != 1200 {
		t.Fatalf("%v %v", e.ActualWidth, e.ActualHeight)
	}
	if err := e.Validate(); err != nil {
		t.Fatal(err)
	}
	c, err := e.Clone()
	if err != nil {
		t.Fatal(err)
	}
	if len(c.(*ApertureEntity).ActualWidth) != 2 {
		t.Fatal("clone lost width")
	}
	e.Quantity = 3
	if e.Validate() == nil {
		t.Fatal("expected size mismatch error")
	}
	if len(e.GetUpdates(c)) == 0 {
		t.Fatal("expected quantity update")
	}
}
