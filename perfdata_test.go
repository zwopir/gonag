package gonag

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

// 'label'=value[UOM];[warn];[crit];[min];[max]
var perfdataTestTable = []struct {
	perfdataItem string
	perfdata     *PerfdataItem
}{
	{"label=3;;;;", &PerfdataItem{
		Label:      "label",
		Value:      "3",
		Thresholds: nil,
		UOM:        &numbersUOM{},
	}},
	{"label=3.0;;;;", &PerfdataItem{
		Label:      "label",
		Value:      "3.0",
		Thresholds: nil,
		UOM:        &numbersUOM{},
	}},
	{"'label with blanks'=3c;;;;", &PerfdataItem{
		Label:      "label with blanks",
		Value:      "3",
		Thresholds: nil,
		UOM:        &countsUOM{},
	}},
	{"n=3;4;5;0;10", &PerfdataItem{
		Label: "n",
		Value: "3",
		Thresholds: Thresholds{
			Warn: "4",
			Crit: "5",
			Min:  "0",
			Max:  "10",
		},
		UOM: &numbersUOM{},
	}},
}

func TestNewPerfdataItem(t *testing.T) {
	for _, tt := range perfdataTestTable {
		actual, err := NewPerfdataItem(tt.perfdataItem)
		if err != nil {
			t.Errorf("parsing perdata item string %s failed with %q",
				tt.perfdataItem, err)
		}
		if actual.Label != tt.perfdata.Label {
			t.Errorf("NewPerfdataItem(%s): expected %s as label, but got %v",
				tt.perfdataItem, tt.perfdata.Label, actual.Label)
		}
		if actual.Value != tt.perfdata.Value {
			t.Errorf("NewPerfdataItem(%s): expected %v as value, but got %v",
				tt.perfdataItem, tt.perfdata.Value, actual.Value)
		}
		for idx, th := range actual.Thresholds {
			if th != tt.perfdata.Thresholds[idx] {
				t.Errorf("Got %s as PerfdataItem %s-Threshold for %q, expected %s",
					th, idx, tt.perfdataItem, tt.perfdata.Thresholds[idx])
			} else {
				t.Logf("Threshold (%s) of %q is %s", idx, tt.perfdataItem, th)
			}
		}
	}
}



func TestPerfdataItem_String(t *testing.T) {
	for _, tt := range perfdataTestTable {
		actual := tt.perfdata.String()
		if actual != tt.perfdataItem {
			t.Errorf("expected %q from perfdataItem.String(), got %q",
			tt.perfdataItem, actual)
		}
	}
}

var splitPerfdataTestTable = []struct {
	in  string
	out []string
}{
	{"'a label'=1.0 b=2c c=3", []string{"'a label'=1.0", "b=2c", "c=3"}},
	{"'label'=1;2;3;4;5 foo=3.14", []string{"'label'=1;2;3;4;5", "foo=3.14"}},
}

var newPerfdataTestTable = []struct {
	in  string
	out []*PerfdataItem
}{
	{"'label'=1;2;3;4;5 foo=3.14", Perfdata{
		{
			Label: "label",
			Value: "1",
			Thresholds: Thresholds{
				Warn: "2",
				Crit: "3",
				Min:  "4",
				Max:  "5",
			},
			UOM: &numbersUOM{},
		},
		{
			Label:      "foo",
			Value:      "3.14",
			Thresholds: Thresholds{
				Warn: "",
				Crit: "",
				Min: "",
				Max: "",
			},
			UOM:        &numbersUOM{},
		},
	}},
}

func TestNewPerfdata(t *testing.T) {
	for _, tt := range newPerfdataTestTable {
		actual, err := NewPerfdata(tt.in)
		if err != nil {
			t.Errorf(err.Error())
		}
		if !reflect.DeepEqual(actual, tt.out) {
			t.Errorf("Expected %s parsing perfdata string, got %s",
				actual, tt.out)
		}
	}

}

func TestGetPerfdataSplitFunc(t *testing.T) {
	splitFunc := GetPerfdataSplitFunc(" ")
	for _, tt := range splitPerfdataTestTable {
		out := []string{}
		scanner := bufio.NewScanner(strings.NewReader(tt.in))
		scanner.Split(splitFunc)
		for scanner.Scan() {
			out = append(out, scanner.Text())
		}
		if !reflect.DeepEqual(out, tt.out) {
			t.Errorf("GetPerfdataSplitFunc failed to split %q. Got %v, expected %v",
				tt.in, out, tt.out)
			for idx, v := range out {
				t.Logf("retrieved element: %d = %q", idx, v)

			}
		}
	}
}

var PerfdataIdentifierTestTable = []struct {
	in  PerfdataThresholdIdentifier
	out string
}{
	{Warn, "Warn"},
	{Crit, "Crit"},
	{Min, "Min"},
	{Max, "Max"},
}

func TestPerfdataThresholdIdentifier_String(t *testing.T) {
	for _, tt := range PerfdataIdentifierTestTable {
		if tt.in.String() != tt.out {
			t.Errorf("PerfdataThresholdIdentifier String() method returned %q, expected %q",
				tt.in.String(), tt.out)
		}
	}
}



