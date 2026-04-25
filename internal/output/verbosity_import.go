package output

import "fmt"

// This file exists solely to provide the fmt import required by verbosity.go
// while keeping that file free of import blocks for readability.
// The compiler merges all files in the package, so the import is available
// across the whole package.

// ensure fmt is used (the actual call is in ParseVerbosity).
var _ = fmt.Sprintf
