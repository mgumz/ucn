PROJECT=ucn
PKG=github.com/mgumz/$(PROJECT)
VERSION=$(shell cat VERSION)
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_HASH=$(shell git describe --tags --always --dirty --match="v*")

TARGETS=linux.amd64 	\
	linux.386 			\
	linux.arm64 		\
	linux.mips64 		\
	darwin.amd64 		\
	darwin.arm64		\
	windows.amd64.exe 	\
	freebsd.amd64

BINARIES=$(addprefix bin/$(PROJECT)-$(VERSION)., $(TARGETS))
RELEASES=$(subst windows.amd64.tar.gz,windows.amd64.zip,$(foreach r,$(subst .exe,,$(TARGETS)),releases/$(PROJECT)-$(VERSION).$(r).tar.gz))
RELEASES += releases/ucn-$(VERSION).alfredworkflow

LDFLAGS=-trimpath -ldflags "-s -w -X $(PKG)/internal/pkg.Version=$(VERSION) \
	-X $(PKG)/internal/pkg.BuildDate=$(BUILD_DATE) \
	-X $(PKG)/internal/pkg.GitHash=$(GIT_HASH)"

$(PROJECT): bin/$(PROJECT)

######################################################
## release related

binaries: $(BINARIES)
release: $(RELEASES)
releases: $(RELEASES)
list-releases:
	@echo $(RELEASES)|tr ' ' '\n'
clean:
	rm -f $(BINARIES) $(RELEASES)

bin/$(PROJECT): cmd/$(PROJECT) bin
	go build $(LDFLAGS) -v -o $@ ./$<

bin/$(PROJECT)-$(VERSION).%:
	env GOARCH=$(subst .,,$(suffix $(subst .exe,,$@))) GOOS=$(subst .,,$(suffix $(basename $(subst .exe,,$@)))) CGO_ENABLED=0 \
	go build $(LDFLAGS) -o $@ ./cmd/$(PROJECT)

releases/$(PROJECT)-$(VERSION).%.zip: bin/$(PROJECT)-$(VERSION).%.exe
	mkdir -p releases
	zip -9 -j -r $@ README.md LICENSE $<
releases/$(PROJECT)-$(VERSION).%.tar.gz: bin/$(PROJECT)-$(VERSION).%
	mkdir -p releases
	tar -cf $(basename $@) README.md LICENSE && \
		tar -rf $(basename $@) --strip-components 1 $< && \
		gzip -9 $(basename $@)

bin:
	mkdir $@

generate-internals:
	go generate -v ./internal/...

generate-cmd:
	go generate -v ./cmd/...

######################################################
## dev related

deps-vendor:
	go mod vendor
deps-cleanup:
	go mod tidy
deps-ls:
	go list -m -mod=readonly -f '{{if not .Indirect}}{{.}}{{end}}' all
deps-ls-updates:
	go list -m -mod=readonly -f '{{if not .Indirect}}{{.}}{{end}}' -u all



compile-analysis: cmd/$(PROJECT)
	go build -gcflags '-m' ./$^

# https://github.com/nektos/act
run-github-workflow-lint:
	act -j lint --container-architecture linux/amd64
run-github-workflow-test:
	act -j test --container-architecture linux/amd64
run-github-workflow-buildLinux:
	act -j buildLinux --container-architecture linux/amd64

reports: report-golangci-lint
reports: report-vuln report-vet

report-golangci-lint:
	@echo '####################################################################'
	golangci-lint run ./cmd/... ./internal/...

report-vuln:
	@echo '####################################################################'
	govulncheck ./cmd/... ./internal/...

report-grype:
	@echo '####################################################################'
	grype .

fetch-report-tools:
	go install golang.org/x/vuln/cmd/govulncheck@latest

fetch-report-tool-grype:
	go install github.com/anchore/grype@latest


test:
	go test -v ./internal/... ./cmd/...

.PHONY: $(PROJECT) bin/$(PROJECT) binaries releases

######################################################
## extra related

bin/ucn.darwin.amd64: bin/ucn-$(VERSION).darwin.amd64
	cp $< $@
bin/ucn.darwin.arm64: bin/ucn-$(VERSION).darwin.arm64
	cp $< $@


ALFRED_FILES=info.plist sshot-fs8.png icon.png

releases/ucn-$(VERSION).alfredworkflow:
releases/ucn-$(VERSION).alfredworkflow: bin/ucn.darwin.amd64
releases/ucn-$(VERSION).alfredworkflow: bin/ucn.darwin.arm64
releases/ucn-$(VERSION).alfredworkflow:
	mkdir -p releases
	rm -f $@
	( cd extra/ucn.alfredworkflow && zip -v -9 -r ../../$@ $(ALFRED_FILES) )
	( cd bin && zip -v -9 ../$@ ucn.darwin.amd64 ucn.darwin.arm64 )
