package validator

import "golang.org/x/mod/semver"

func ValidateVersion(version string) bool {
	return semver.IsValid(version)
}
