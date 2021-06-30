package models

import (
	"reflect"
	"testing"
	"time"
)

func Test_FetchMeta(t *testing.T) {
	var tests = []struct {
		in       FetchMeta
		outdated bool
	}{
		{
			in: FetchMeta{
				SchemaVersion: 0,
			},
			outdated: true,
		},
		{
			in: FetchMeta{
				SchemaVersion: LatestSchemaVersion,
			},
			outdated: false,
		},
	}

	for i, tt := range tests {
		if aout := tt.in.OutDated(); tt.outdated != aout {
			t.Errorf("[%d] outdated expected: %#v\n  actual: %#v\n", i, tt.outdated, aout)
		}
	}
}

func Test_ConvertOSV(t *testing.T) {
	var tests = []struct {
		in       []OSVJSON
		expected []OSV
	}{
		{
			in: []OSVJSON{
				{
					ID:        "RUSTSEC-2016-0001",
					Published: "2016-11-05T12:00:00Z",
					Modified:  "2020-10-02T01:29:11Z",
					// Withdrawn: "",
					Aliases: []string{"CVE-2016-10931"},
					Related: []string{},
					Package: struct {
						Ecosystem string "json:\"ecosystem\""
						Name      string "json:\"name\""
						Purl      string "json:\"purl\""
					}{
						Name:      "openssl",
						Ecosystem: "crates.io",
						Purl:      "pkg:cargo/openssl"},
					Summary: "Summary",
					Details: "Details",
					Affects: struct {
						Ranges []struct {
							Type       string "json:\"type\""
							Repo       string "json:\"repo\""
							Introduced string "json:\"introduced\""
							Fixed      string "json:\"fixed\""
						} "json:\"ranges\""
						Versions []string "json:\"versions\""
					}{
						Ranges: []struct {
							Type       string "json:\"type\""
							Repo       string "json:\"repo\""
							Introduced string "json:\"introduced\""
							Fixed      string "json:\"fixed\""
						}{
							{
								Type:  "SEMVER",
								Fixed: "0.9.0",
							},
						},
						Versions: []string{},
					},
					References: []struct {
						Type string "json:\"type\""
						URL  string "json:\"url\""
					}{

						{
							Type: "PACKAGE",
							URL:  "https://crates.io/crates/openssl",
						},
						{
							Type: "ADVISORY",
							URL:  "https://rustsec.org/advisories/RUSTSEC-2016-0001.html",
						},
					},
					// Severity: "",
					EcosystemSpecific: map[string]map[string][]string{
						"affects": {
							"functions": []string{},
							"arch":      []string{},
							"os":        []string{},
						},
					},
					DatabaseSpecific: map[string]interface{}{
						"informational": nil,
						"source":        "https://github.com/Shnatsel/advisory-db/blob/osv/crates/RUSTSEC-2016-0001.json",
						"cvss":          nil,
						"categories":    []string{},
					},
				},
			},
			expected: []OSV{
				{
					EntryID:   "RUSTSEC-2016-0001",
					Published: time.Date(2016, 11, 5, 12, 0, 0, 0, time.UTC),
					Modified:  time.Date(2020, 10, 2, 1, 29, 11, 0, time.UTC),
					Withdrawn: time.Date(1000, time.January, 1, 0, 0, 0, 0, time.UTC),
					Aliases: []OSVAliases{
						{Alias: "CVE-2016-10931"},
					},
					Related: []OSVRelated{},
					Package: OSVPackage{
						Name:      "openssl",
						Ecosystem: "crates.io",
						Purl:      "pkg:cargo/openssl",
					},
					Summary: "Summary",
					Details: "Details",
					Affects: OSVAffects{
						Ranges: []OSVAffectsRanges{
							{
								Type:  "SEMVER",
								Fixed: "0.9.0",
							},
						},
						Versions: []OSVAffectsVersions{},
					},
					References: []OSVReferences{
						{
							Type: "PACKAGE",
							URL:  "https://crates.io/crates/openssl",
						}, {
							Type: "ADVISORY",
							URL:  "https://rustsec.org/advisories/RUSTSEC-2016-0001.html",
						},
					},
					// Severity: "",
					EcosystemSpecific: OSVEcosystemSpecific{},
					DatabaseSpecific:  OSVDatabaseSpecific{},
				},
			},
		},
	}

	for i, tt := range tests {
		aout, err := ConvertOSV(tt.in)
		if err != nil {
			t.Errorf("[%d] ConvertOSV error: %w", i, err)
		}

		if !reflect.DeepEqual(tt.expected, aout) {
			t.Errorf("[%d] ConvertOSV expected: %#v\n  actual: %#v\n", i, tt.expected, aout)
		}
	}
}
