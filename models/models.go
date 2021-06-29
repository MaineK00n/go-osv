package models

import "gorm.io/gorm"

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
