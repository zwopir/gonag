package gonag

import (
	"testing"
)

var uomTT = []struct {
	in         string
	value      string
	unitstring string
	baseunit   string
	magnitude  int
}{
	{"1.23", "1.23", "", "", 0},
	{"1", "1", "", "", 0},
	{"2B", "2", "B", "B", 0},
	{"2.34c", "2.34", "c", "c", 0},
	{"2c", "2", "c", "c", 0},
	{"3.45B", "3.45", "B", "B", 0},
	{"3.45KB", "3.45", "KB", "B", 3},
	{"3.45MB", "3.45", "MB", "B", 6},
	{"3.45GB", "3.45", "GB", "B", 9},
	{"3.45TB", "3.45", "TB", "B", 12},
	{"4.56%", "4.56", "%", "%", 0},
	{"5.67s", "5.67", "s", "s", 0},
	{"5.67ms", "5.67", "ms", "s", -3},
	{"5.67us", "5.67", "us", "s", -6},
}

var failingUOMTT = []string{
	"not a number",
	"1a",
	"1,2",
	"1..2",
	"ms",
}

func TestParseValue(t *testing.T) {
	for _, tt := range uomTT {
		value, uom, err := ParseValue(tt.in)
		if err != nil {
			t.Errorf("parsing %q returned an error: %s", tt.in, err)
		}
		if value != tt.value {
			t.Errorf("parsing %q failed. Expected %s as value, got %s", tt.in, tt.value, value)
		}
		if uom.String() != tt.unitstring {
			t.Errorf("parsing %q failed. Expected %s as unitstring, got %s", tt.in, tt.unitstring, uom.String())
		}
		if uom.Baseunit() != tt.baseunit {
			t.Errorf("parsing %q failed. Expected %s as baseunit, got %s", tt.in, tt.baseunit, uom.Baseunit())
		}
		if uom.Magnitude() != tt.magnitude {
			t.Errorf("parsing %q failed. Expected %s as magnitude, got %s", tt.in, tt.magnitude, uom.Magnitude())
		}
	}
	for _, inputString := range failingUOMTT {
		value, uom, err := ParseValue(inputString)
		if err == nil {
			t.Errorf("input string %q should raise an error", inputString)
			t.Logf("String() returns %s", uom.String())
			t.Logf("Baseunit() returns %s", uom.Baseunit())
			t.Logf("Magnitude() returns %s", uom.Magnitude())
			t.Logf("value returns %s", value)
		} else {
			t.Logf("correctly received an error parsing %q: %s", inputString, err)
		}
	}
}
