package output

// Verbosity controls how much detail is emitted during a drift-check run.
type Verbosity int

const (
	// VerbositySilent suppresses all non-error output.
	VerbositySilent Verbosity = iota
	// VerbosityNormal emits the summary and any diff.
	VerbosityNormal
	// VerbosityVerbose emits per-stage progress in addition to the summary.
	VerbosityVerbose
	// VerbosityDebug emits internal timing and raw manifests as well.
	VerbosityDebug
)

var verbosityNames = map[Verbosity]string{
	VerbositySilent:  "silent",
	VerbosityNormal:  "normal",
	VerbosityVerbose: "verbose",
	VerbosityDebug:   "debug",
}

// String returns the human-readable name of the verbosity level.
func (v Verbosity) String() string {
	if name, ok := verbosityNames[v]; ok {
		return name
	}
	return "unknown"
}

// ParseVerbosity converts a string such as "verbose" into the corresponding
// Verbosity constant. It returns VerbosityNormal and an error when the string
// is not recognised.
func ParseVerbosity(s string) (Verbosity, error) {
	for v, name := range verbosityNames {
		if name == s {
			return v, nil
		}
	}
	return VerbosityNormal, fmt.Errorf("unknown verbosity %q: must be one of silent, normal, verbose, debug", s)
}

// KnownVerbosities returns every recognised verbosity name in a stable order.
func KnownVerbosities() []string {
	return []string{
		verbosityNames[VerbositySilent],
		verbosityNames[VerbosityNormal],
		verbosityNames[VerbosityVerbose],
		verbosityNames[VerbosityDebug],
	}
}
