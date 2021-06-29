package models

import (
	"gorm.io/gorm"
)

// LatestSchemaVersion manages the Schema version used in the latest Go-OSV.
const LatestSchemaVersion = 1

// FetchMeta has meta infomation about fetched security tracker
type FetchMeta struct {
	gorm.Model    `json:"-"`
	GoOSVRevision string
	SchemaVersion uint
}

// OutDated checks whether last fetched feed is out dated
func (f FetchMeta) OutDated() bool {
	return f.SchemaVersion != LatestSchemaVersion
}

// OSVType :
type OSVType string

const (
	// CratesIOType :
	CratesIOType OSVType = "crates.io"
	// DWFType :
	DWFType OSVType = "DWF"
	// GoType :
	GoType OSVType = "Go"
	// LinuxType :
	LinuxType OSVType = "Linux"
	// OSSFuzzType :
	OSSFuzzType OSVType = "OSS-Fuzz"
	// PyPIType :
	PyPIType OSVType = "PyPI"
)

// OSVJSON : https://osv.dev/docs/#tag/vulnerability_schema
type OSVJSON struct {
	ID        string   `json:"id"`
	Published string   `json:"published"`
	Modified  string   `json:"modified"`
	Withdrawn string   `json:"withdrawn"`
	Aliases   []string `json:"aliases"`
	Related   []string `json:"related"`
	Package   struct {
		Ecosystem string `json:"ecosystem"`
		Name      string `json:"name"`
		Purl      string `json:"purl"`
	} `json:"package"`
	Summary string `json:"summary"`
	Details string `json:"details"`
	Affects struct {
		Ranges []struct {
			Type       string `json:"type"`
			Repo       string `json:"repo"`
			Introduced string `json:"introduced"`
			Fixed      string `json:"fixed"`
		} `json:"ranges"`
		Versions []string `json:"versions"`
	} `json:"affects"`
	References []struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"references"`
	Severity          string      `json:"severity"`
	EcosystemSpecific interface{} `json:"ecosystem_specific"`
	DatabaseSpecific  interface{} `json:"database_specific"`
}
