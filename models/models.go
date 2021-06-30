package models

import (
	"time"

	"golang.org/x/xerrors"
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

// OSV :
type OSV struct {
	ID                int64  `json:"-" gorm:"index:idx_osv_id"`
	EntryID           string `json:"ID"`
	Published         time.Time
	Modified          time.Time
	Withdrawn         time.Time
	Aliases           []OSVAliases
	Related           []OSVRelated
	Package           OSVPackage
	Summary           string
	Details           string
	Affects           OSVAffects
	References        []OSVReferences
	Severity          string
	EcosystemSpecific OSVEcosystemSpecific
	DatabaseSpecific  OSVDatabaseSpecific
}

// OSVAliases :
type OSVAliases struct {
	ID    int64  `json:"-"`
	OSVID int64  `json:"-"`
	Alias string `gorm:"index:idx_osv_aliases_alias"`
}

// OSVRelated :
type OSVRelated struct {
	ID      int64 `json:"-"`
	OSVID   int64 `json:"-"`
	Related string
}

// OSVPackage :
type OSVPackage struct {
	ID        int64  `json:"-"`
	OSVID     int64  `json:"-"`
	Ecosystem string `gorm:"index:idx_osv_packages_ecosystem"`
	Name      string `gorm:"index:idx_osv_packages_name"`
	Purl      string
}

// OSVAffects :
type OSVAffects struct {
	ID       int64 `json:"-"`
	OSVID    int64 `json:"-"`
	Ranges   []OSVAffectsRanges
	Versions []OSVAffectsVersions
}

// OSVAffectsRanges :
type OSVAffectsRanges struct {
	ID           int64 `json:"-"`
	OSVAffectsID int64 `json:"-"`
	Type         string
	Repo         string
	Introduced   string
	Fixed        string
}

// OSVAffectsVersions :
type OSVAffectsVersions struct {
	ID           int64 `json:"-"`
	OSVAffectsID int64 `json:"-"`
	Version      string
}

// OSVReferences :
type OSVReferences struct {
	ID    int64 `json:"-"`
	OSVID int64 `json:"-"`
	Type  string
	URL   string
}

// OSVEcosystemSpecific :
type OSVEcosystemSpecific struct {
	ID    int64 `json:"-"`
	OSVID int64 `json:"-"`
}

// OSVDatabaseSpecific :
type OSVDatabaseSpecific struct {
	ID    int64 `json:"-"`
	OSVID int64 `json:"-"`
}

// ConvertOSV :
func ConvertOSV(osvJSONs []OSVJSON) (osvs []OSV, err error) {
	for _, osvJSON := range osvJSONs {
		osv := OSV{
			EntryID: osvJSON.ID,
			Package: OSVPackage{
				Ecosystem: osvJSON.Package.Ecosystem,
				Name:      osvJSON.Package.Name,
				Purl:      osvJSON.Package.Purl,
			},
			Summary:  osvJSON.Summary,
			Details:  osvJSON.Details,
			Severity: osvJSON.Severity,
		}

		if osv.Published, err = time.Parse(time.RFC3339, osvJSON.Published); err != nil {
			return []OSV{}, xerrors.Errorf("Failed to time.Parse. err: %w", err)
		}

		if osv.Modified, err = time.Parse(time.RFC3339, osvJSON.Modified); err != nil {
			return []OSV{}, xerrors.Errorf("Failed to time.Parse. err: %w", err)
		}

		if osvJSON.Withdrawn != "" {
			if osv.Withdrawn, err = time.Parse(time.RFC3339, osvJSON.Withdrawn); err != nil {
				return []OSV{}, xerrors.Errorf("Failed to time.Parse. err: %w", err)
			}
		} else {
			osv.Withdrawn = time.Date(1000, time.January, 1, 0, 0, 0, 0, time.UTC)
		}

		aliases := []OSVAliases{}
		for _, a := range osvJSON.Aliases {
			aliases = append(aliases, OSVAliases{Alias: a})
		}
		osv.Aliases = aliases

		related := []OSVRelated{}
		for _, r := range osvJSON.Related {
			related = append(related, OSVRelated{Related: r})
		}
		osv.Related = related

		affects := OSVAffects{}

		ranges := []OSVAffectsRanges{}
		for _, r := range osvJSON.Affects.Ranges {
			ranges = append(ranges, OSVAffectsRanges{
				Type:       r.Type,
				Repo:       r.Repo,
				Introduced: r.Introduced,
				Fixed:      r.Fixed,
			})
		}
		affects.Ranges = ranges

		versions := []OSVAffectsVersions{}
		for _, v := range osvJSON.Affects.Versions {
			versions = append(versions, OSVAffectsVersions{Version: v})
		}
		affects.Versions = versions

		osv.Affects = affects

		references := []OSVReferences{}
		for _, r := range osvJSON.References {
			references = append(references, OSVReferences{
				Type: r.Type,
				URL:  r.URL,
			})
		}
		osv.References = references

		osvs = append(osvs, osv)
	}

	return osvs, nil
}
