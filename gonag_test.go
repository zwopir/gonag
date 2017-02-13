package gonag

import (
	"testing"
)

var ReturnCodeStringerTestTable = []struct{
	in ReturnCode
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


