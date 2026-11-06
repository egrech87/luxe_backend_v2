package int64Array

import (
	"reflect"
	"testing"
)

func TestScanValue(t *testing.T) {
	var a Int64Array
	if err := a.Scan("{900,905}"); err != nil || !reflect.DeepEqual(a, Int64Array{900, 905}) {
		t.Fatalf("scan: %v %v", a, err)
	}
	if err := a.Scan("{}"); err != nil || len(a) != 0 {
		t.Fatalf("empty scan: %v %v", a, err)
	}
	v, _ := Int64Array{1, 2}.Value()
	if v != "{1,2}" {
		t.Fatalf("value: %v", v)
	}
	if v, _ := Int64Array(nil).Value(); v != "{}" {
		t.Fatalf("nil value: %v", v)
	}
}

func TestCollapsedAndArea(t *testing.T) {
	if got := (Int64Array{900, 900}).Collapsed("/"); got != "900" {
		t.Fatal(got)
	}
	if got := (Int64Array{900, 905}).Collapsed("/"); got != "900/905" {
		t.Fatal(got)
	}
	if got := AreaSqm(Int64Array{1000, 2000}, Int64Array{1000, 1000}); got != 3 {
		t.Fatal(got)
	}
}
