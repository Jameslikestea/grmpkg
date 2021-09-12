package model

import "time"

type Package struct {
	Name     string      `json:"name"`
	Info     PackageInfo `json:"info"`
	Versions VersionList `json:"versions"`
}

type PackageInfo struct {
	Package  string `json:"package"`
	Hostname string `json:"hostname"`
}

type PackageVersion struct {
	Name    string `json:"Name"`
	Short   string `json:"Short"`
	Version string `json:"Version"`
	Time    string `json:"Time"`
}

type VersionList []PackageVersion

func (vl VersionList) Len() int {
	return len(vl)
}

func (vl VersionList) Less(i, j int) bool {
	date1, _ := time.Parse(time.RFC3339, vl[i].Time)
	date2, _ := time.Parse(time.RFC3339, vl[j].Time)

	return date1.After(date2)
}

func (vl VersionList) Swap(i, j int) {
	vl[i], vl[j] = vl[j], vl[i]
}
