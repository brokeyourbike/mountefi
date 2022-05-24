// The version package provides a location to set the release versions for all
// packages to consume, without creating import cycles.
//
// This package should not import any other packages.
package version

import (
	"github.com/hashicorp/go-version"
)

// The main version number that is being run at the moment.
var Version = "0.0.0"

// SemVer is an instance of version.Version. This has the secondary
// benefit of verifying during tests and init time that our version is a
// proper semantic version, which should always be the case.
var SemVer *version.Version

func init() {
	SemVer = version.Must(version.NewVersion(Version))
}

// Header is the header name used to send the current mountefi version
// in http requests.
const Header = "Mountefi-Version"

// String returns the complete version string
func String() string {
	return Version
}
