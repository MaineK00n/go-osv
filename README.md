# go-osv
`go-osv` builds a local copy of [Open Source Vulnerabilities; OSV](https://osv.dev/list).

# Abstract
`go-osv` is written in Go, and therefore you can just grab the binary releases and drop it in your $PATH.

`go-osv` builds a local copy of [Open Source Vulnerabilities; OSV](https://osv.dev/list).

# Main features
`go-osv` has the following features.
- Build a local copy of Open Source Vulnerabilities; OSV
- A server mode for easy querying

# Installation
## Requirements
- SQLite3, MySQL, PostgreSQL or Redis
- git
- go 
## Install
```console
$ mkdir -p $GOPATH/src/github.com/MaineK00n
$ cd $GOPATH/src/github.com/MaineK00n
$ git clone https://github.com/MaineK00n/go-osv.git
$ cd go-osv
$ make install
```

# Usage
```console
$ go-osv
Open Source Vulnerabilities;OSV

Usage:
  go-osv [command]

Available Commands:
  fetch       Fetch the data of the osv-vulnerabilities
  help        Help about any command
  server      Start OSV HTTP server

Flags:
      --config string       config file (default is $HOME/.go-osv.yaml)
      --dbpath string       /path/to/sqlite3 or SQL connection string (default "/home/mainek00n/github/github.com/MaineK00n/go-osv/go-osv.sqlite3")
      --dbtype string       Database type to store data in (sqlite3, mysql, postgres or redis supported) (default "sqlite3")
      --debug               debug mode (default: false)
      --debug-sql           SQL debug mode
  -h, --help                help for go-osv
      --http-proxy string   http://proxy-url:port (default: empty)
      --log-dir string      /path/to/log (default "/var/log/go-osv")
      --log-json            output log as JSON

Use "go-osv [command] --help" for more information about a command.
```

# Fetch osv-vulnerabilities/crates.io
```console
$ go-osv fetch crates.io
INFO[06-30|11:44:27] Initialize Database 
INFO[06-30|11:44:27] Fetched all OSV Data from osv-vulnerabilities/crates.io 
INFO[06-30|11:44:28] Fetched                                  OSVs=289
INFO[06-30|11:44:28] Insert OSVs into DB                      db=sqlite3
 289 / 289 [========================================================] 100.00% 0s
```

# Fetch osv-vulnerabilities/DWF
```console
$ go-osv fetch dwf
INFO[06-30|11:44:44] Initialize Database 
INFO[06-30|11:44:44] Fetched all OSV Data from osv-vulnerabilities/DWF 
INFO[06-30|11:44:45] Fetched                                  OSVs=15
INFO[06-30|11:44:45] Insert OSVs into DB                      db=sqlite3
 15 / 15 [==========================================================] 100.00% 0s

```

# Fetch osv-vulnerabilities/Go
```console
$ go-osv fetch go
INFO[06-30|11:44:59] Initialize Database 
INFO[06-30|11:44:59] Fetched all OSV Data from osv-vulnerabilities/Go 
INFO[06-30|11:44:59] Fetched                                  OSVs=92
INFO[06-30|11:44:59] Insert OSVs into DB                      db=sqlite3
 92 / 92 [==========================================================] 100.00% 0s
```

# Fetch osv-vulnerabilities/Linux
```console
$ go-osv fetch linux
INFO[06-30|11:45:12] Initialize Database 
INFO[06-30|11:45:12] Fetched all OSV Data from osv-vulnerabilities/Linux 
INFO[06-30|11:45:12] Fetched                                  OSVs=811
INFO[06-30|11:45:12] Insert OSVs into DB                      db=sqlite3
 811 / 811 [========================================================] 100.00% 0s
```

# Fetch osv-vulnerabilities/OSS-Fuzz
```console
$ go-osv fetch oss-fuzz
INFO[06-30|11:45:28] Initialize Database 
INFO[06-30|11:45:28] Fetched all OSV Data from osv-vulnerabilities/OSS-Fuzz 
INFO[06-30|11:45:29] Fetched                                  OSVs=1592
INFO[06-30|11:45:29] Insert OSVs into DB                      db=sqlite3
 1592 / 1592 [======================================================] 100.00% 0s
```

# Fetch osv-vulnerabilities/PyPI
```console
$ go-osv fetch pypi
INFO[06-30|11:45:42] Initialize Database 
INFO[06-30|11:45:42] Fetched all OSV Data from osv-vulnerabilities/PyPI 
INFO[06-30|11:45:42] Fetched                                  OSVs=451
INFO[06-30|11:45:42] Insert OSVs into DB                      db=sqlite3
 451 / 451 [========================================================] 100.00% 0s
```


# Server mode
```console
$ go-osv server
INFO[06-30|11:46:03] Starting HTTP Server... 
INFO[06-30|11:46:03] Listening                                URL=127.0.0.1:1328

   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v3.3.10-dev
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on 127.0.0.1:1328

$ curl http://127.0.0.1:1328/ids/CVE-2016-10931 | jq
[
  {
    "ID": "RUSTSEC-2016-0001",
    "Published": "2016-11-05T12:00:00Z",
    "Modified": "2020-10-02T01:29:11Z",
    "Withdrawn": "1000-01-01T00:00:00Z",
    "Aliases": [
      {
        "Alias": "CVE-2016-10931"
      }
    ],
    "Related": [],
    "Package": {
      "Ecosystem": "crates.io",
      "Name": "openssl",
      "Purl": "pkg:cargo/openssl"
    },
    "Summary": "SSL/TLS MitM vulnerability due to insecure defaults",
    "Details": "All versions of rust-openssl prior to 0.9.0 contained numerous insecure defaults\nincluding off-by-default certificate verification and no API to perform hostname\nverification.\n\nUnless configured correctly by a developer, these defaults could allow an attacker\nto perform man-in-the-middle attacks.\n\nThe problem was addressed in newer versions by enabling certificate verification\nby default and exposing APIs to perform hostname verification. Use the\n`SslConnector` and `SslAcceptor` types to take advantage of these new features\n(as opposed to the lower-level `SslContext` type).",
    "Affects": {
      "Ranges": [
        {
          "Type": "SEMVER",
          "Repo": "",
          "Introduced": "",
          "Fixed": "0.9.0"
        }
      ],
      "Versions": []
    },
    "References": [
      {
        "Type": "PACKAGE",
        "URL": "https://crates.io/crates/openssl"
      },
      {
        "Type": "ADVISORY",
        "URL": "https://rustsec.org/advisories/RUSTSEC-2016-0001.html"
      },
      {
        "Type": "WEB",
        "URL": "https://github.com/sfackler/rust-openssl/releases/tag/v0.9.0"
      }
    ],
    "Severity": "",
    "EcosystemSpecific": {},
    "DatabaseSpecific": {}
  }
]

$ curl http://127.0.0.1:1328/crates.io/ids/CVE-2016-10931 | jq
[
  {
    "ID": "RUSTSEC-2016-0001",
    "Published": "2016-11-05T12:00:00Z",
    "Modified": "2020-10-02T01:29:11Z",
    "Withdrawn": "1000-01-01T00:00:00Z",
    "Aliases": [
      {
        "Alias": "CVE-2016-10931"
      }
    ],
    "Related": [],
    "Package": {
      "Ecosystem": "crates.io",
      "Name": "openssl",
      "Purl": "pkg:cargo/openssl"
    },
    "Summary": "SSL/TLS MitM vulnerability due to insecure defaults",
    "Details": "All versions of rust-openssl prior to 0.9.0 contained numerous insecure defaults\nincluding off-by-default certificate verification and no API to perform hostname\nverification.\n\nUnless configured correctly by a developer, these defaults could allow an attacker\nto perform man-in-the-middle attacks.\n\nThe problem was addressed in newer versions by enabling certificate verification\nby default and exposing APIs to perform hostname verification. Use the\n`SslConnector` and `SslAcceptor` types to take advantage of these new features\n(as opposed to the lower-level `SslContext` type).",
    "Affects": {
      "Ranges": [
        {
          "Type": "SEMVER",
          "Repo": "",
          "Introduced": "",
          "Fixed": "0.9.0"
        }
      ],
      "Versions": []
    },
    "References": [
      {
        "Type": "PACKAGE",
        "URL": "https://crates.io/crates/openssl"
      },
      {
        "Type": "ADVISORY",
        "URL": "https://rustsec.org/advisories/RUSTSEC-2016-0001.html"
      },
      {
        "Type": "WEB",
        "URL": "https://github.com/sfackler/rust-openssl/releases/tag/v0.9.0"
      }
    ],
    "Severity": "",
    "EcosystemSpecific": {},
    "DatabaseSpecific": {}
  }
]

$ curl http://127.0.0.1:1328/pkgs/openssl | jq
[
  {
    "ID": "RUSTSEC-2016-0001",
    "Published": "2016-11-05T12:00:00Z",
    "Modified": "2020-10-02T01:29:11Z",
    "Withdrawn": "1000-01-01T00:00:00Z",
    "Aliases": [
      {
        "Alias": "CVE-2016-10931"
      }
    ],
    "Related": [],
    "Package": {
      "Ecosystem": "crates.io",
      "Name": "openssl",
      "Purl": "pkg:cargo/openssl"
    },
    "Summary": "SSL/TLS MitM vulnerability due to insecure defaults",
    "Details": "All versions of rust-openssl prior to 0.9.0 contained numerous insecure defaults\nincluding off-by-default certificate verification and no API to perform hostname\nverification.\n\nUnless configured correctly by a developer, these defaults could allow an attacker\nto perform man-in-the-middle attacks.\n\nThe problem was addressed in newer versions by enabling certificate verification\nby default and exposing APIs to perform hostname verification. Use the\n`SslConnector` and `SslAcceptor` types to take advantage of these new features\n(as opposed to the lower-level `SslContext` type).",
    "Affects": {
      "Ranges": [
        {
          "Type": "SEMVER",
          "Repo": "",
          "Introduced": "",
          "Fixed": "0.9.0"
        }
      ],
      "Versions": []
    },
    "References": [
      {
        "Type": "PACKAGE",
        "URL": "https://crates.io/crates/openssl"
      },
      {
        "Type": "ADVISORY",
        "URL": "https://rustsec.org/advisories/RUSTSEC-2016-0001.html"
      },
      {
        "Type": "WEB",
        "URL": "https://github.com/sfackler/rust-openssl/releases/tag/v0.9.0"
      }
    ],
    "Severity": "",
    "EcosystemSpecific": {},
    "DatabaseSpecific": {}
  },
  {
    "ID": "RUSTSEC-2018-0010",
    "Published": "2018-06-01T12:00:00Z",
    "Modified": "2020-10-02T01:29:11Z",
    "Withdrawn": "1000-01-01T00:00:00Z",
    "Aliases": [
      {
        "Alias": "CVE-2018-20997"
      }
    ],
    "Related": [],
    "Package": {
      "Ecosystem": "crates.io",
      "Name": "openssl",
      "Purl": "pkg:cargo/openssl"
    },
    "Summary": "Use after free in CMS Signing",
    "Details": "Affected versions of the OpenSSL crate used structures after they'd been freed.",
    "Affects": {
      "Ranges": [
        {
          "Type": "SEMVER",
          "Repo": "",
          "Introduced": "0.10.8",
          "Fixed": "0.10.9"
        }
      ],
      "Versions": []
    },
    "References": [
      {
        "Type": "PACKAGE",
        "URL": "https://crates.io/crates/openssl"
      },
      {
        "Type": "ADVISORY",
        "URL": "https://rustsec.org/advisories/RUSTSEC-2018-0010.html"
      },
      {
        "Type": "WEB",
        "URL": "https://github.com/sfackler/rust-openssl/pull/942"
      }
    ],
    "Severity": "",
    "EcosystemSpecific": {},
    "DatabaseSpecific": {}
  },
  {
    "ID": "OSV-2018-109",
    "Published": "2021-01-13T00:00:48.206043Z",
    "Modified": "2021-03-09T04:49:04.82093Z",
    "Withdrawn": "1000-01-01T00:00:00Z",
    "Aliases": [],
    "Related": [],
    "Package": {
      "Ecosystem": "OSS-Fuzz",
      "Name": "openssl",
      "Purl": ""
    },
    "Summary": "Heap-use-after-free in ssl_get_prev_session",
    "Details": "OSS-Fuzz report: https://bugs.chromium.org/p/oss-fuzz/issues/detail?id=8241\n\nCrash type: Heap-use-after-free READ 4\nCrash state:\nssl_get_prev_session\ntls_early_post_process_client_hello\ntls_post_process_client_hello\n",
    "Affects": {
      "Ranges": [
        {
          "Type": "GIT",
          "Repo": "https://github.com/openssl/openssl.git",
          "Introduced": "61fb59238dad6452a37ec14513fae617a4faef29",
          "Fixed": "5f96a95e2562f026557f625e50c052e77c7bc2e8"
        }
      ],
      "Versions": []
    },
    "References": [
      {
        "Type": "REPORT",
        "URL": "https://bugs.chromium.org/p/oss-fuzz/issues/detail?id=8241"
      }
    ],
    "Severity": "",
    "EcosystemSpecific": {},
    "DatabaseSpecific": {}
  },
  {
    "ID": "OSV-2018-153",
    "Published": "2021-01-13T00:01:05.75724Z",
    "Modified": "2021-06-23T06:28:23.524218Z",
    "Withdrawn": "1000-01-01T00:00:00Z",
    "Aliases": [],
    "Related": [],
    "Package": {
      "Ecosystem": "OSS-Fuzz",
      "Name": "openssl",
      "Purl": ""
    },
    "Summary": "Heap-buffer-overflow in asn1_ex_i2c",
    "Details": "OSS-Fuzz report: https://bugs.chromium.org/p/oss-fuzz/issues/detail?id=7696\n\nCrash type: Heap-buffer-overflow READ 4\nCrash state:\nasn1_ex_i2c\nasn1_i2d_ex_primitive\nASN1_item_ex_i2d\n",
    "Affects": {
      "Ranges": [
        {
          "Type": "GIT",
          "Repo": "https://github.com/openssl/openssl.git",
          "Introduced": "902f7d5c87d66a78d3eb10709c6cb3486a216b48",
          "Fixed": "0df65d82dbc41e8da00adb243de5918db532c8a6"
        }
      ],
      "Versions": [
        {
          "Version": "OpenSSL_1_1_1-pre1"
        },
        {
          "Version": "OpenSSL_1_1_1-pre2"
        },
        {
          "Version": "OpenSSL_1_1_1-pre3"
        },
        {
          "Version": "OpenSSL_1_1_1-pre4"
        },
        {
          "Version": "OpenSSL_1_1_1-pre5"
        },
        {
          "Version": "OpenSSL_1_1_1-pre6"
        },
        {
          "Version": "OpenSSL_1_1_1-pre7"
        }
      ]
    },
    "References": [
      {
        "Type": "REPORT",
        "URL": "https://bugs.chromium.org/p/oss-fuzz/issues/detail?id=7696"
      }
    ],
    "Severity": "",
    "EcosystemSpecific": {},
    "DatabaseSpecific": {}
  },
  {
    "ID": "OSV-2020-223",
    "Published": "2020-06-24T01:51:19.666966Z",
    "Modified": "2021-03-09T04:49:05.731028Z",
    "Withdrawn": "1000-01-01T00:00:00Z",
    "Aliases": [],
    "Related": [],
    "Package": {
      "Ecosystem": "OSS-Fuzz",
      "Name": "openssl",
      "Purl": ""
    },
    "Summary": "Heap-use-after-free in CRYPTO_DOWN_REF",
    "Details": "OSS-Fuzz report: https://bugs.chromium.org/p/oss-fuzz/issues/detail?id=21550\n\nCrash type: Heap-use-after-free WRITE 4\nCrash state:\nCRYPTO_DOWN_REF\nDH_free\nevp_pkey_free_legacy\n",
    "Affects": {
      "Ranges": [
        {
          "Type": "GIT",
          "Repo": "https://github.com/openssl/openssl.git",
          "Introduced": "ada66e78ef535fe80e422bbbadffe8e7863d457c",
          "Fixed": "fe56d5951f0b42fd3ff1cf42a96d07f06f9692bc"
        }
      ],
      "Versions": []
    },
    "References": [
      {
        "Type": "REPORT",
        "URL": "https://bugs.chromium.org/p/oss-fuzz/issues/detail?id=21550"
      }
    ],
    "Severity": "",
    "EcosystemSpecific": {},
    "DatabaseSpecific": {}
  },
  {
    "ID": "OSV-2020-29",
    "Published": "2020-06-24T01:51:10.908381Z",
    "Modified": "2021-03-09T04:49:05.807418Z",
    "Withdrawn": "1000-01-01T00:00:00Z",
    "Aliases": [],
    "Related": [],
    "Package": {
      "Ecosystem": "OSS-Fuzz",
      "Name": "openssl",
      "Purl": ""
    },
    "Summary": "Heap-use-after-free in CRYPTO_DOWN_REF",
    "Details": "OSS-Fuzz report: https://bugs.chromium.org/p/oss-fuzz/issues/detail?id=20816\n\nCrash type: Heap-use-after-free WRITE 4\nCrash state:\nCRYPTO_DOWN_REF\nDH_free\nevp_pkey_free_it\n",
    "Affects": {
      "Ranges": [
        {
          "Type": "GIT",
          "Repo": "https://github.com/openssl/openssl.git",
          "Introduced": "ada66e78ef535fe80e422bbbadffe8e7863d457c",
          "Fixed": "fe56d5951f0b42fd3ff1cf42a96d07f06f9692bc"
        }
      ],
      "Versions": []
    },
    "References": [
      {
        "Type": "REPORT",
        "URL": "https://bugs.chromium.org/p/oss-fuzz/issues/detail?id=20816"
      }
    ],
    "Severity": "",
    "EcosystemSpecific": {},
    "DatabaseSpecific": {}
  },
  {
    "ID": "OSV-2020-386",
    "Published": "2020-07-01T00:00:06.528477Z",
    "Modified": "2021-03-09T04:49:05.859492Z",
    "Withdrawn": "1000-01-01T00:00:00Z",
    "Aliases": [],
    "Related": [],
    "Package": {
      "Ecosystem": "OSS-Fuzz",
      "Name": "openssl",
      "Purl": ""
    },
    "Summary": "Heap-buffer-overflow in OPENSSL_strlcpy",
    "Details": "OSS-Fuzz report: https://bugs.chromium.org/p/oss-fuzz/issues/detail?id=16107\n\nCrash type: Heap-buffer-overflow WRITE 1\nCrash state:\nOPENSSL_strlcpy\nOPENSSL_strlcat\nERR_add_error_vdata\n",
    "Affects": {
      "Ranges": [
        {
          "Type": "GIT",
          "Repo": "https://github.com/openssl/openssl.git",
          "Introduced": "10f8b36874fca928c3f41834babac8ee94dd3f09",
          "Fixed": "036913b1076da41f257c640a5e6230476c647eff"
        }
      ],
      "Versions": []
    },
    "References": [
      {
        "Type": "REPORT",
        "URL": "https://bugs.chromium.org/p/oss-fuzz/issues/detail?id=16107"
      }
    ],
    "Severity": "",
    "EcosystemSpecific": {},
    "DatabaseSpecific": {}
  },
  {
    "ID": "OSV-2020-430",
    "Published": "2020-07-01T00:00:09.096641Z",
    "Modified": "2021-03-09T04:49:05.883624Z",
    "Withdrawn": "1000-01-01T00:00:00Z",
    "Aliases": [],
    "Related": [],
    "Package": {
      "Ecosystem": "OSS-Fuzz",
      "Name": "openssl",
      "Purl": ""
    },
    "Summary": "Stack-use-after-return in OSSL_PARAM_get_int32",
    "Details": "OSS-Fuzz report: https://bugs.chromium.org/p/oss-fuzz/issues/detail?id=15114\n\nCrash type: Stack-use-after-return READ 4\nCrash state:\nOSSL_PARAM_get_int32\nmd5_sha1_set_params\nssl3_final_finish_mac\n",
    "Affects": {
      "Ranges": [
        {
          "Type": "GIT",
          "Repo": "https://github.com/openssl/openssl.git",
          "Introduced": "d5e5e2ffafc7dbc861f7d285508cf129c5e8f5ac",
          "Fixed": "83b4a24384e62ed8cf91f51bf9a303f98017e13e"
        }
      ],
      "Versions": []
    },
    "References": [
      {
        "Type": "REPORT",
        "URL": "https://bugs.chromium.org/p/oss-fuzz/issues/detail?id=15114"
      }
    ],
    "Severity": "",
    "EcosystemSpecific": {},
    "DatabaseSpecific": {}
  },
  {
    "ID": "OSV-2020-442",
    "Published": "2020-07-01T00:00:09.812508Z",
    "Modified": "2021-03-09T04:49:05.89008Z",
    "Withdrawn": "1000-01-01T00:00:00Z",
    "Aliases": [],
    "Related": [],
    "Package": {
      "Ecosystem": "OSS-Fuzz",
      "Name": "openssl",
      "Purl": ""
    },
    "Summary": "Heap-buffer-overflow in CRYPTO_strdup",
    "Details": "OSS-Fuzz report: https://bugs.chromium.org/p/oss-fuzz/issues/detail?id=17715\n\nCrash type: Heap-buffer-overflow READ 14\nCrash state:\nCRYPTO_strdup\nX509V3_add_value\ni2v_GENERAL_NAME\n",
    "Affects": {
      "Ranges": [
        {
          "Type": "GIT",
          "Repo": "https://github.com/openssl/openssl.git",
          "Introduced": "5053a3766a13f40afb3c89f54d1f9a5eae38a3eb",
          "Fixed": "aec9667bd19a8ca9bdd519db3a231a95b9e92674"
        }
      ],
      "Versions": []
    },
    "References": [
      {
        "Type": "REPORT",
        "URL": "https://bugs.chromium.org/p/oss-fuzz/issues/detail?id=17715"
      }
    ],
    "Severity": "",
    "EcosystemSpecific": {},
    "DatabaseSpecific": {}
  }
]

$ curl http://127.0.0.1:1328/crates.io/pkgs/openssl | jq
[
  {
    "ID": "RUSTSEC-2016-0001",
    "Published": "2016-11-05T12:00:00Z",
    "Modified": "2020-10-02T01:29:11Z",
    "Withdrawn": "1000-01-01T00:00:00Z",
    "Aliases": [
      {
        "Alias": "CVE-2016-10931"
      }
    ],
    "Related": [],
    "Package": {
      "Ecosystem": "crates.io",
      "Name": "openssl",
      "Purl": "pkg:cargo/openssl"
    },
    "Summary": "SSL/TLS MitM vulnerability due to insecure defaults",
    "Details": "All versions of rust-openssl prior to 0.9.0 contained numerous insecure defaults\nincluding off-by-default certificate verification and no API to perform hostname\nverification.\n\nUnless configured correctly by a developer, these defaults could allow an attacker\nto perform man-in-the-middle attacks.\n\nThe problem was addressed in newer versions by enabling certificate verification\nby default and exposing APIs to perform hostname verification. Use the\n`SslConnector` and `SslAcceptor` types to take advantage of these new features\n(as opposed to the lower-level `SslContext` type).",
    "Affects": {
      "Ranges": [
        {
          "Type": "SEMVER",
          "Repo": "",
          "Introduced": "",
          "Fixed": "0.9.0"
        }
      ],
      "Versions": []
    },
    "References": [
      {
        "Type": "PACKAGE",
        "URL": "https://crates.io/crates/openssl"
      },
      {
        "Type": "ADVISORY",
        "URL": "https://rustsec.org/advisories/RUSTSEC-2016-0001.html"
      },
      {
        "Type": "WEB",
        "URL": "https://github.com/sfackler/rust-openssl/releases/tag/v0.9.0"
      }
    ],
    "Severity": "",
    "EcosystemSpecific": {},
    "DatabaseSpecific": {}
  },
  {
    "ID": "RUSTSEC-2018-0010",
    "Published": "2018-06-01T12:00:00Z",
    "Modified": "2020-10-02T01:29:11Z",
    "Withdrawn": "1000-01-01T00:00:00Z",
    "Aliases": [
      {
        "Alias": "CVE-2018-20997"
      }
    ],
    "Related": [],
    "Package": {
      "Ecosystem": "crates.io",
      "Name": "openssl",
      "Purl": "pkg:cargo/openssl"
    },
    "Summary": "Use after free in CMS Signing",
    "Details": "Affected versions of the OpenSSL crate used structures after they'd been freed.",
    "Affects": {
      "Ranges": [
        {
          "Type": "SEMVER",
          "Repo": "",
          "Introduced": "0.10.8",
          "Fixed": "0.10.9"
        }
      ],
      "Versions": []
    },
    "References": [
      {
        "Type": "PACKAGE",
        "URL": "https://crates.io/crates/openssl"
      },
      {
        "Type": "ADVISORY",
        "URL": "https://rustsec.org/advisories/RUSTSEC-2018-0010.html"
      },
      {
        "Type": "WEB",
        "URL": "https://github.com/sfackler/rust-openssl/pull/942"
      }
    ],
    "Severity": "",
    "EcosystemSpecific": {},
    "DatabaseSpecific": {}
  }
]
```

# Contribute

1. fork a repository: github.com/MaineK00n/go-osv to github.com/you/repo
2. get original code: `go get github.com/MaineK00n/go-osv`
3. work on original code
4. add remote to your repo: git remote add myfork https://github.com/you/repo.git
5. push your changes: git push myfork
6. create a new Pull Request

- see [GitHub and Go: forking, pull requests, and go-getting](http://blog.campoy.cat/2014/03/github-and-go-forking-pull-requests-and.html)

----

# License
MIT

# Author
[MaineK00n](https://twitter.com/MaineK00n)