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

generate-docs:
	go generate ./cmd/$(PROJECT)

generate-internals:
	go generate ./internal/...

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

reports: report-vuln report-gosec
reports: report-staticcheck report-vet report-ineffassign
reports: report-cyclo
reports: report-errcheck report-gocritic
reports: report-misspell

report-cyclo:
	@echo '####################################################################'
	gocyclo ./cmd/$(PROJECT) ./internal/pkg
report-misspell:
	@echo '####################################################################'
	misspell ./cmd/... ./internal/...
report-ineffassign:
	@echo '####################################################################'
	ineffassign ./cmd/... ./internal/...
report-vet:
	@echo '####################################################################'
	go vet ./cmd/... ./internal/...
report-staticcheck:
	@echo '####################################################################'
	staticcheck ./cmd/... ./internal/...
report-vuln:
	@echo '####################################################################'
	govulncheck ./cmd/... ./internal/...
report-gosec:
	@echo '####################################################################'
	gosec ./cmd/... ./internal/...
report-grype:
	@echo '####################################################################'
	grype .
report-errcheck:
	@echo '####################################################################'
	errcheck -ignorepkg fmt ./...
report-gocritic:
	@echo '####################################################################'
	gocritic check ./cmd/... ./internal/...

fetch-report-tools:
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	go install github.com/client9/misspell/cmd/misspell@latest
	go install github.com/gordonklaus/ineffassign@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go install -v github.com/go-critic/go-critic/cmd/gocritic@latest
	go install github.com/kisielk/errcheck@latest

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

releases/ucn-$(VERSION).alfredworkflow: bin/ucn.darwin.amd64
releases/ucn-$(VERSION).alfredworkflow: bin/ucn.darwin.arm64
releases/ucn-$(VERSION).alfredworkflow:
	rm -f $@
	( cd extra/ucn.alfredworkflow && zip -v -9 -r ../../$@ . )
	( cd bin && zip -v -9 ../$@ ucn.darwin.amd64 ucn.darwin.arm64 )
