package gonag

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

/*
func NewFromPluginOutput(returnCode ReturnCode, pluginOutput string) (*CheckResult, error){
	parts := strings.SplitAfterN(pluginOutput, "|", 2)
	perfdata, err := NewPerfdataItem(parts[1])
	return &CheckResult{
		Text: parts[0],
		ReturnCode: returnCode,
		Perfdata: perfdata,
	}, err
}
*/
