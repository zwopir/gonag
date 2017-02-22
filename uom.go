package gonag

import (
	"fmt"
	"strconv"
	"strings"
)

type numbersUOM struct{}

func (*numbersUOM) String() string   { return "" }
func (*numbersUOM) Baseunit() string { return "" }
func (*numbersUOM) Magnitude() int   { return 0 }

type countsUOM struct{}

func (*countsUOM) String() string   { return "c" }
func (*countsUOM) Baseunit() string { return "c" }
func (*countsUOM) Magnitude() int   { return 0 }

type bytesUOM struct {
	magnitude int
}

func (b *bytesUOM) String() string {
	switch b.magnitude {
	case 3:
		return "KB"
	case 6:
		return "MB"
	case 9:
		return "GB"
	case 12:
		return "TB"
	}
	return "B"
}

func (*bytesUOM) Baseunit() string { return "B" }

func (b *bytesUOM) Magnitude() int {
	return b.magnitude
}

type secondsUOM struct {
	magnitude int
}

func (s *secondsUOM) String() string {
	switch s.magnitude {
	case -6:
		return "us"
	case -3:
		return "ms"
	}
	return "s"
}

func (*secondsUOM) Baseunit() string { return "s" }

func (s *secondsUOM) Magnitude() int {
	return s.magnitude
}

type percentUOM struct{}

func (*percentUOM) String() string   { return "%" }
func (*percentUOM) Baseunit() string { return "%" }
func (*percentUOM) Magnitude() int   { return 0 }

func parseUnitString(unitString string) (Uniter, error) {
	switch {
	// no UOM given, return a base uniter
	case len(unitString) == 0:
		return &numbersUOM{}, nil
	case len(unitString) == 1:
		switch unitString {
		case "s":
			return &secondsUOM{magnitude: 0}, nil
		case "B":
			return &bytesUOM{magnitude: 0}, nil
		case "c":
			return &countsUOM{}, nil
		case "%":
			return &percentUOM{}, nil
		}
		return nil, fmt.Errorf("unknown single character UOM string %s", unitString)
	case unitString == "ms":
		return &secondsUOM{magnitude: -3}, nil
	case unitString == "us":
		return &secondsUOM{magnitude: -6}, nil
	case unitString == "KB":
		return &bytesUOM{magnitude: 3}, nil
	case unitString == "MB":
		return &bytesUOM{magnitude: 6}, nil
	case unitString == "GB":
		return &bytesUOM{magnitude: 9}, nil
	case unitString == "TB":
		return &bytesUOM{magnitude: 12}, nil
	}
	return nil, fmt.Errorf("unknown UOM string %s", unitString)
}

func ParseValue(s string) (string, Uniter, error) {
	unitString := ""
	if strings.Contains("csB%", string(s[len(s)-1])) {
		unitString = string(s[len(s)-1])
		if strings.Contains("muKMGT", string(s[len(s)-2:len(s)-1])) {
			unitString = string(s[len(s)-2:])
		}
	}
	valueString := strings.TrimSuffix(s, unitString)
	_, err := strconv.ParseFloat(valueString, 32)
	if err != nil {
		return "", nil, err
	}
	uom, err := parseUnitString(unitString)
	if err != nil {
		return "", nil, err
	}
	return valueString, uom, nil
}


