// Package buildinfo captures version information injected at link time.
package buildinfo

import (
	"fmt"
	"runtime"
)

// These variables are set via -ldflags at build time.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// Info describes the built binary.
type Info struct {
	Version   string
	Commit    string
	Date      string
	GoVersion string
	Platform  string
}

// Get returns the build information for this binary.
func Get() Info {
	return Info{
		Version:   version,
		Commit:    commit,
		Date:      date,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// String renders the build info as a single line.
func (i Info) String() string {
	return fmt.Sprintf("build: %s commit=%s date=%s %s %s",
		i.Version, i.Commit, i.Date, i.GoVersion, i.Platform)
}
