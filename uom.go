package gonag

import (
	"fmt"
	"strings"
	"strconv"
)

type numbers struct{}
func (*numbers) String() string {return ""}
func (*numbers) Baseunit() string {return ""}
func (*numbers) Magnitude() int {return 0}

type counts struct{}
func (*counts) String() string {return "c"}
func (*counts) Baseunit() string {return "c"}
func (*counts) Magnitude() int {return 0}

type bytes struct{
	magnitude int
}

func (b *bytes) String() string {
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

func (*bytes) Baseunit() string {return "B"}

func (b *bytes) Magnitude() int {
	return b.magnitude
}

type seconds struct{
	magnitude int
}

func (s *seconds) String() string {
	switch s.magnitude {
	case -6:
		return "us"
	case -3:
		return "ms"
	}
	return "s"
}

func (*seconds) Baseunit() string {return "s"}

func (s *seconds) Magnitude() int {
	return s.magnitude
}

type percent struct{}
func (*percent) String() string {return "%"}
func (*percent) Baseunit() string {return "%"}
func (*percent) Magnitude() int {return 0}

func parseUnitString(unitString string) (Uniter, error) {
	switch {
	// no UOM given, return a base uniter
	case len(unitString) == 0:
		return &numbers{}, nil
	case len(unitString) == 1:
		switch unitString {
		case "s":
			return &seconds{magnitude:0}, nil
		case "B":
			return &bytes{magnitude:0}, nil
		case "c":
			return &counts{}, nil
		case "%":
			return &percent{}, nil
		}
		return nil, fmt.Errorf("unknown single character UOM string %s", unitString)
	case unitString == "ms":
		return &seconds{magnitude:-3}, nil
	case unitString == "us":
		return &seconds{magnitude:-6}, nil
	case unitString == "KB":
		return &bytes{magnitude:3}, nil
	case unitString == "MB":
		return &bytes{magnitude:6}, nil
	case unitString == "GB":
		return &bytes{magnitude:9}, nil
	case unitString == "TB":
		return &bytes{magnitude:12}, nil
	}
	return nil, fmt.Errorf("unknown UOM string %s", unitString)
}



func ParseValue(s string) (string, Uniter, error) {
	unitString := ""
	if strings.Contains("csB%", string(s[len(s)-1])) {
		unitString = string(s[len(s) - 1])
		if strings.Contains("muKMGT", string(s[len(s) - 2:len(s) - 1])) {
			unitString = string(s[len(s) - 2:])
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

