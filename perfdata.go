package gonag

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type PerfdataThresholdIdentifier int

//go:generate stringer -type=PerfdataThresholdIdentifier
const (
	Warn PerfdataThresholdIdentifier = iota
	Crit
	Min
	Max
)

var (
	ErrUnknownValue = errors.New("Unknown perfdata value")
	ErrEmptyValue   = errors.New("Empty perfdata value")
)

type Thresholds map[PerfdataThresholdIdentifier]string

// Perfdata represents the performance data of a Nagios check
type Perfdata []*PerfdataItem

func (pd Perfdata) String() string {
	out := []string{}
	for _, s := range pd {
		out = append(out, fmt.Sprint(s))
	}
	return strings.Join(out, " ")
}

// PerfdataItem represents a single performance data item of a Nagios check. Values (Value, Warn, Crit, Min, Max) are encoded as
// string, since there can be U (unknown), a number or not set (which is not the default initialization of a float/int)
type PerfdataItem struct {
	Label      string
	Value      string
	Thresholds Thresholds
	UOM        Uniter
}

func (pd *PerfdataItem) String() string {
	label := pd.Label
	if strings.Contains(label, " ") {
		label = fmt.Sprintf("'%s'", label)
	}
	return fmt.Sprintf("%s=%s%s;%s;%s;%s;%s",
		label, pd.Value, pd.UOM,
		pd.Thresholds[Warn],
		pd.Thresholds[Crit],
		pd.Thresholds[Min],
		pd.Thresholds[Max],
	)
}

func NewPerfdataItem(perfdataItem string) (*PerfdataItem, error) {
	perfdata := PerfdataItem{
		Thresholds: Thresholds{},
	}
	// initialize empty thresholds
	for idx := 0; idx <= 3; idx++ {
		perfdata.Thresholds[PerfdataThresholdIdentifier(idx)] = ""
	}
	parts := strings.SplitAfterN(perfdataItem, "=", 2)
	perfdata.Label = strings.Trim(parts[0], "'=")
	values := strings.SplitAfter(parts[1], ";")
	if values[0] == "" {
		return &perfdata, nil
	} else {
		value, uom, err := ParseValue(strings.TrimSuffix(values[0], ";"))
		if err != nil {
			return nil, fmt.Errorf("Error parseing value with UOM: %s", err)
		}
		perfdata.UOM = uom
		perfdata.Value = value
	}

	for idx, value := range values[1:] {
		if value == "" {
			continue
		}
		v := strings.TrimSuffix(value, ";")
		// ToDo: is silently ignoring parse errors ok here?
		if _, err := strconv.ParseFloat(v, 32); err == nil {
			perfdata.Thresholds[PerfdataThresholdIdentifier(idx)] = v
		}
	}
	return &perfdata, nil
}

func NewPerfdata(perfdataString string) ([]*PerfdataItem, error) {
	perfdata := []*PerfdataItem{}
	splitFunc := GetPerfdataSplitFunc(" ")
	scanner := bufio.NewScanner(strings.NewReader(perfdataString))
	scanner.Split(splitFunc)
	for scanner.Scan() {
		perfdataItem, err := NewPerfdataItem(scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("Failed to parse perfdata item %q: %s",
				scanner.Text(), err)
		}
		perfdata = append(perfdata, perfdataItem)
	}
	return perfdata, nil
}

func GetPerfdataSplitFunc(splitter string) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {

		var accumulatedData []byte
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		indexOfEqualSign := strings.Index(string(data), "=")
		if indexOfEqualSign >= 0 {
			accumulatedData = append(accumulatedData, data[0:indexOfEqualSign]...)
		} else {
			return 0, nil, fmt.Errorf("error parsing perfdata items")
		}
		endOfPerfdataItem := strings.Index(string(data[indexOfEqualSign:]), splitter)
		if endOfPerfdataItem >= 0 {
			accumulatedData = append(
				accumulatedData, data[indexOfEqualSign:indexOfEqualSign+endOfPerfdataItem]...,
			)
			return indexOfEqualSign + endOfPerfdataItem + 1, accumulatedData, nil
		}

		if atEOF {
			return len(data), data, nil
		}

		return
	}
}
