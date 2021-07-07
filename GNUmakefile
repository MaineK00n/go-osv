.PHONY: \
	build \
	install \
	all \
	vendor \
	lint \
	vet \
	fmt \
	fmtcheck \
	pretest \
	test \
	integration \
	cov \
	clean \
	build-integration \
	clean-integration \
	fetch-rdb \
	fetch-redis \
	diff-cveid \
	diff-package \
	diff-server-rdb \
	diff-server-redis \
	diff-server-rdb-redis

SRCS = $(shell git ls-files '*.go')
PKGS = $(shell go list ./...)
VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'github.com/MaineK00n/go-osv/config.Version=$(VERSION)' \
	-X 'github.com/MaineK00n/go-osv/config.Revision=$(REVISION)'
GO := GO111MODULE=on go
GO_OFF := GO111MODULE=off go

all: build

build: main.go pretest
	$(GO) build -ldflags "$(LDFLAGS)" -o go-osv  $<

install: main.go pretest
	$(GO) install -ldflags "$(LDFLAGS)"

b: 	main.go pretest
	$(GO) build -ldflags "$(LDFLAGS)" -o go-osv $<

lint:
	$(GO_OFF) get -u golang.org/x/lint/golint
	golint $(PKGS)

vet:
	echo $(PKGS) | xargs env $(GO) vet || exit;

fmt:
	gofmt -s -w $(SRCS)

mlint:
	$(foreach file,$(SRCS),gometalinter $(file) || exit;)

fmtcheck:
	$(foreach file,$(SRCS),gofmt -s -d $(file);)

pretest: lint vet fmtcheck

test: 
	$(GO) test -cover -v ./... || exit;

unused:
	$(foreach pkg,$(PKGS),unused $(pkg);)

integration:
	go test -tags docker_integration -run TestIntegration -v

cov:
	@ go get -v github.com/axw/gocov/gocov
	@ go get golang.org/x/tools/cmd/cover
	gocov test | gocov report

clean:
	$(foreach pkg,$(PKGS),go clean $(pkg) || exit;)

BRANCH := $(shell git symbolic-ref --short HEAD)
build-integration:
	@ git stash save
	$(GO) build -ldflags "$(LDFLAGS)" -o integration/go-osv.new
	git checkout $(shell git describe --tags --abbrev=0)
	@git reset --hard
	$(GO) build -ldflags "$(LDFLAGS)" -o integration/go-osv.old
	git checkout $(BRANCH)
	-@ git stash apply stash@{0} && git stash drop stash@{0}

clean-integration:
	-pkill go-osv.old
	-pkill go-osv.new
	-rm integration/go-osv.old integration/go-osv.new integration/go-osv.old.sqlite3 integration/go-osv.new.sqlite3
	-docker kill redis-old redis-new
	-docker rm redis-old redis-new

fetch-rdb:
	integration/go-osv.old fetch crates.io --dbpath=integration/go-osv.old.sqlite3
	integration/go-osv.old fetch dwf --dbpath=integration/go-osv.old.sqlite3
	integration/go-osv.old fetch go --dbpath=integration/go-osv.old.sqlite3
	integration/go-osv.old fetch linux --dbpath=integration/go-osv.old.sqlite3
	integration/go-osv.old fetch oss-fuzz --dbpath=integration/go-osv.old.sqlite3
	integration/go-osv.old fetch pypi --dbpath=integration/go-osv.old.sqlite3

	integration/go-osv.new fetch crates.io --dbpath=integration/go-osv.new.sqlite3
	integration/go-osv.new fetch dwf --dbpath=integration/go-osv.new.sqlite3
	integration/go-osv.new fetch go --dbpath=integration/go-osv.new.sqlite3
	integration/go-osv.new fetch linux --dbpath=integration/go-osv.new.sqlite3
	integration/go-osv.new fetch oss-fuzz --dbpath=integration/go-osv.new.sqlite3
	integration/go-osv.new fetch pypi --dbpath=integration/go-osv.new.sqlite3

fetch-redis:
	docker run --name redis-old -d -p 127.0.0.1:6379:6379 redis
	docker run --name redis-new -d -p 127.0.0.1:6380:6379 redis

	integration/go-osv.old fetch crates.io --dbtype redis --dbpath "redis://127.0.0.1:6379/0"
	integration/go-osv.old fetch dwf --dbtype redis --dbpath "redis://127.0.0.1:6379/0"
	integration/go-osv.old fetch go --dbtype redis --dbpath "redis://127.0.0.1:6379/0"
	integration/go-osv.old fetch linux --dbtype redis --dbpath "redis://127.0.0.1:6379/0"
	integration/go-osv.old fetch oss-fuzz --dbtype redis --dbpath "redis://127.0.0.1:6379/0"
	integration/go-osv.old fetch pypi --dbtype redis --dbpath "redis://127.0.0.1:6379/0"

	integration/go-osv.new fetch crates.io --dbtype redis --dbpath "redis://127.0.0.1:6380/0"
	integration/go-osv.new fetch dwf --dbtype redis --dbpath "redis://127.0.0.1:6380/0"
	integration/go-osv.new fetch go --dbtype redis --dbpath "redis://127.0.0.1:6380/0"
	integration/go-osv.new fetch linux --dbtype redis --dbpath "redis://127.0.0.1:6380/0"
	integration/go-osv.new fetch oss-fuzz --dbtype redis --dbpath "redis://127.0.0.1:6380/0"
	integration/go-osv.new fetch pypi --dbtype redis --dbpath "redis://127.0.0.1:6380/0"

diff-id:
	@ python integration/diff_server_mode.py id all
	@ python integration/diff_server_mode.py id crates.io
	@ python integration/diff_server_mode.py id DWF
	@ python integration/diff_server_mode.py id Go
	@ python integration/diff_server_mode.py id Linux
	@ python integration/diff_server_mode.py id OSS-Fuzz
	@ python integration/diff_server_mode.py id PyPI


diff-package:
	@ python integration/diff_server_mode.py package all
	@ python integration/diff_server_mode.py package crates.io
	@ python integration/diff_server_mode.py package DWF
	@ python integration/diff_server_mode.py package Go
	@ python integration/diff_server_mode.py package Linux
	@ python integration/diff_server_mode.py package OSS-Fuzz
	@ python integration/diff_server_mode.py package PyPI

diff-server-rdb:
	integration/go-osv.old server --dbpath=integration/go-osv.old.sqlite3 --port 1325 > /dev/null & 
	integration/go-osv.new server --dbpath=integration/go-osv.new.sqlite3 --port 1326 > /dev/null &
	make diff-id
	make diff-package
	pkill go-osv.old 
	pkill go-osv.new

diff-server-redis:
	integration/go-osv.old server --dbtype redis --dbpath "redis://127.0.0.1:6379/0" --port 1325 > /dev/null & 
	integration/go-osv.new server --dbtype redis --dbpath "redis://127.0.0.1:6380/0" --port 1326 > /dev/null &
	make diff-id
	make diff-package
	pkill go-osv.old 
	pkill go-osv.new

diff-server-rdb-redis:
	integration/go-osv.new server --dbpath=integration/go-osv.new.sqlite3 --port 1325 > /dev/null &
	integration/go-osv.new server --dbtype redis --dbpath "redis://127.0.0.1:6380/0" --port 1326 > /dev/null &
	make diff-id
	make diff-package
	pkill go-osv.new