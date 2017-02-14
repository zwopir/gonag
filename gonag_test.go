package gonag

import (
	"testing"
	"reflect"
)

var ReturnCodeStringerTestTable = []struct {
	in  ReturnCode
	out string
}{
	{OK, "OK"},
	{WARNING, "WARNING"},
	{CRITICAL, "CRITICAL"},
	{UNKNOWN, "UNKNOWN"},
}

func TestReturnCode_String(t *testing.T) {
	for _, tt := range ReturnCodeStringerTestTable {
		if tt.in.String() != tt.out {
			t.Errorf("Returncode String() method returned %q, expected %q",
				tt.in.String(), tt.out)
		}
	}
}

var pluginOutputTestTable = []struct {
	pluginOutput string
	exitcode ReturnCode
	expected *CheckResult
}{
	{"plugin text with blanks|a=123.3c n=5;4;6;0;10 free=8MB", OK, &CheckResult{
		Text: "plugin text with blanks",
		ReturnCode: OK,
		Perfdata: []*Perfdata{
			{
				Label: "a",
				Value: "123.3",
				Thresholds: Thresholds{},
				UOM: &counts{},
			},
			{
				Label: "n",
				Value: "5",
				Thresholds: Thresholds{
					Warn: "4",
					Crit: "6",
					Min: "0",
					Max: "10",
				},
				UOM: &numbers{},
			},
			{
				Label: "free",
				Value: "8",
				Thresholds: Thresholds{},
				UOM: &bytes{magnitude: 3},
			},
		},
	}},
}

func TestNewFromPluginOutput(t *testing.T) {
	for _, tt := range pluginOutputTestTable {
		actual, err := NewFromPluginOutput(tt.exitcode, tt.pluginOutput)
		if err != nil {
			t.Errorf("Parsing plugin output failed: %s", err)
		}
		if actual.ReturnCode != tt.expected.ReturnCode {
			t.Errorf("Returncode is %s, expected %s", actual.ReturnCode, tt.expected.ReturnCode)
		}
		if actual.Text != tt.expected.Text {
			t.Errorf("Text is %q, expected %q", actual.Text, tt.expected.Text)
		}
		if reflect.DeepEqual(tt.expected.Perfdata, actual.Perfdata) {
			t.Errorf("Perfdata is %v, expected %v", actual.Perfdata, tt.expected.Perfdata)
		}
	}
}