package gonag

import (
	"strings"
	"text/template"
	"os"
)

type Uniter interface {
	String() string
	Baseunit() string
	Magnitude() int
}

// ReturnCode represents a Nagios return code
// ReturnCode implements Stringer via go:generate stringer
type ReturnCode int

//go:generate stringer -type=ReturnCode
const (
	OK ReturnCode = iota
	WARNING
	CRITICAL
	UNKNOWN
)

// CheckResult represents a Nagios Check Result
type CheckResult struct {
	Text       string
	ReturnCode ReturnCode
	Perfdata   []*Perfdata
}


func NewFromPluginOutput(returnCode ReturnCode, pluginOutput string) (*CheckResult, error){
	parts := strings.SplitAfterN(pluginOutput, "|", 2)
	text := strings.TrimSuffix(parts[0], "|")
	perfdata, err := NewPerfdata(parts[1])
	if err != nil {
		return nil, err
	}
	return &CheckResult{
		Text: text,
		ReturnCode: returnCode,
		Perfdata: perfdata,
	}, err
}

func (cr *CheckResult) RenderCheckResult(formatString string) (string, error) {
	tmpl, err := template.New("checkResult").Parse(formatString)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(os.Stdout, *cr)
	return "not implemented yet", nil
}