VERSION=$(strip $(subst v, , $(shell git describe --tags --abbrev=0)))
HASH=$(shell git rev-parse --short HEAD)
TIME=$(shell git log --pretty=format:"%ad" --date=format:'%Y%m%d%H%M%S' $(HASH) -1)
BUILD=$(VERSION).$(HASH).$(TIME)
BUILDTIME=$(shell date "+%Y-%m-%dT%H:%M:%S%z")

GO_PACKAGE=$(shell grep "^module" go.mod | awk '{print $$2}')

GO_BUILD_LDFLAGS = -X '$(GO_PACKAGE)/pkg/buildinfo.version=$(VERSION)' \
  				   -X '$(GO_PACKAGE)/pkg/buildinfo.build=$(BUILD)' \
  				   -X '$(GO_PACKAGE)/pkg/buildinfo.buildTime=$(BUILDTIME)'
GO_BUILD_FLAGS	= -ldflags "-w -s $(GO_BUILD_LDFLAGS)"

OS = darwin linux windows
ARCH = arm64 amd64

.PHONY: generate
generate:
	cd cmd/kragins/internal/cli; wire .

.PHONY: build
build: generate
	go build -o "_debug/kragins/" $(GO_PACKAGE)/cmd/kragins;

.PHONY: prod
prod: generate
	for d in $$(go list -f '{{ if (eq .Name "main") }} {{.ImportPath}}  {{end}}' ./cmd/...); do \
		app=$$(basename $${d}); \
		rm -rf dist/$$app; \
		for arch in $(ARCH); do \
  			for os in $(OS); do \
				echo "release [$$app] [$$os]-[$$arch] production"; \
				GOOS=$$os \
				GOARCH=$$arch \
				go build -a -tags prod -o "dist/$$app/$$os/$$arch/" $(GO_BUILD_FLAGS) \
					$(GO_PACKAGE)/cmd/$$app; \
				cd dist/$$app/$$os/$$arch/; \
				zip -r $$app-$$os-$$arch-v$(VERSION).zip . ; \
				cd -; \
			done \
		done \
	done
